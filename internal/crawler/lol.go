package crawler

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"sync"

	"LoLItemRecommender/internal/database"
	"LoLItemRecommender/internal/printer"
	"LoLItemRecommender/internal/queue"
	"LoLItemRecommender/internal/riotapi/api"
	"LoLItemRecommender/internal/riotapi/gamedata"
)

type GameData struct {
	em             *api.EndpointsManager
	client         *api.Client
	staticData     *gamedata.StaticData
	playersCrawled *sync.Map
	db             *database.DB
	lookingChampID int
	//playersData map[string]*gamedata.Player
}

func NewGameData(region string, db *database.DB) (*GameData, error) {
	gd := &GameData{
		em: api.NewEndpointsManager(os.Getenv("RIOT_API_KEY"), region),
		//playersData: make(map[string]*gamedata.Player),
		client:         api.NewClient(),
		playersCrawled: &sync.Map{},
		db:             db,
	}
	gd.staticData = gamedata.NewStaticData(gd.em, gd.client)
	if err := gd.staticData.RetrieveAPIVersion(); err != nil {
		return nil, err
	}
	if err := gd.staticData.RetrieveChampionsStats(); err != nil {
		return nil, err
	}
	if err := gd.staticData.RetrieveItems(); err != nil {
		return nil, err
	}

	var lookingForChampID, _ = strconv.Atoi(gd.staticData.ChampionsStats["Samira"].Key)
	gd.lookingChampID = lookingForChampID
	return gd, nil
}

func (gd *GameData) FlagPlayers(names []string) {
	for _, n := range names {
		gd.playersCrawled.Store(n, true)
	}
}

func (gd *GameData) RetrieveAdditionalPlayerData(player *gamedata.Player) error {
	b, err := gd.client.Get(gd.em.GetSummonerByName(url.PathEscape(player.SummonerName)))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &player); err != nil {
		return err
	}
	return nil
}

func (gd *GameData) RetrievePlayerGamesId(player *gamedata.Player) ([]string, error) {
	b, err := gd.client.Get(gd.em.GetMatchListURL(player.Puuid))
	if err != nil {
		return nil, err
	}
	var matchIDs []string
	if err = json.Unmarshal(b, &matchIDs); err != nil {
		return nil, err
	}
	printer.Info("Found {-F_MAGENTA,BOLD}%d {-RESET}games for {-F_YELLOW}%s", len(matchIDs), player.SummonerName)
	return matchIDs, nil
}

func (gd *GameData) RetrieveGameInfo(gameID string) (*gamedata.MatchData, error) {
	b, err := gd.client.Get(gd.em.GetMatchInfoURL(gameID))
	if err != nil {
		return nil, err
	}
	var matchData = new(gamedata.MatchData)
	if err = json.Unmarshal(b, matchData); err != nil {
		return nil, err
	}
	return matchData, nil
}

func (gd *GameData) InitWithChallengerPlayers() ([]*gamedata.Player, error) {
	b, err := gd.client.Get(gd.em.GetSummonersByLeague(gamedata.RankedSolo5V5, gamedata.Challenger, gamedata.TierOne, 1))
	if err != nil {
		return nil, err
	}
	var p []gamedata.Player
	if err = json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	var pp []*gamedata.Player
	for i := range p {
		pp = append(pp, &p[i])
	}
	return pp, nil
}

func (gd *GameData) CrawlPlayerData(player *gamedata.Player, pool *queue.Pool) error {
	_, ok := gd.playersCrawled.Load(player.SummonerName)
	if ok {
		return nil
	}
	gd.playersCrawled.Store(player.SummonerName, true)
	if player.Puuid == "" {
		if err := gd.RetrieveAdditionalPlayerData(player); err != nil {
			return err
		}
		printer.Debug("Retrieved additional data for player %s", player.SummonerName)
	}
	gameIds, err := gd.RetrievePlayerGamesId(player)
	if err != nil {
		return err
	}
	for _, g := range gameIds {
		matchdata, err := gd.RetrieveGameInfo(g)
		if err != nil {
			return err
		}
		if matchdata.Info.GameMode != gamedata.GameModeClassic || matchdata.Info.GameType != gamedata.GameTypeRanked {
			continue
		}
		for _, p := range matchdata.Info.Participants {
			if p.Role != gamedata.RoleCarry || p.Lane != gamedata.LaneBottom {
				continue
			}
			if p.ChampionId == gd.lookingChampID {
				printer.Info("{-F_GREEN,BOLD}Saving game")
				if err := gd.db.SaveMatch(matchdata); err != nil {
					return err
				}
			}
			if p.SummonerId != player.SummonerId {
				newPlayer := gamedata.Player{
					SummonerId:    p.SummonerId,
					SummonerName:  p.SummonerName,
					SummonerLevel: p.SummonerLevel,
				}
				pool.Dispatch(func() error {
					return gd.CrawlPlayerData(&newPlayer, pool)
				})
			}
		}
	}
	return nil
}
