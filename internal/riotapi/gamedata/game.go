package gamedata

type Perks struct {
	StatPerks struct {
		Defense int `json:"defense"`
		Flex    int `json:"flex"`
		Offense int `json:"offense"`
	} `json:"statPerks"`
	Styles []struct {
		Description string `json:"description"`
		Selections  []struct {
			Perk int `json:"perk"`
			Var1 int `json:"var1"`
			Var2 int `json:"var2"`
			Var3 int `json:"var3"`
		} `json:"selections"`
		Style int `json:"style"`
	} `json:"styles"`
}
type Participant struct {
	Assists                     int    `json:"assists"`
	BaronKills                  int    `json:"baronKills"`
	BountyLevel                 int    `json:"bountyLevel"`
	ChampLevel                  int    `json:"champLevel"`
	ChampionId                  int    `json:"championId"`
	ChampionName                string `json:"championName"`
	DamageDealtToObjectives     int    `json:"damageDealtToObjectives"`
	Deaths                      int    `json:"deaths"`
	GoldEarned                  int    `json:"goldEarned"`
	GoldSpent                   int    `json:"goldSpent"`
	Item0                       int    `json:"item0"`
	Item1                       int    `json:"item1"`
	Item2                       int    `json:"item2"`
	Item3                       int    `json:"item3"`
	Item4                       int    `json:"item4"`
	Item5                       int    `json:"item5"`
	Item6                       int    `json:"item6"`
	Kills                       int    `json:"kills"`
	Lane                        string `json:"lane"`
	ParticipantId               int    `json:"participantId"`
	Perks                       Perks  `json:"perks"`
	Role                        string `json:"role"`
	Summoner1Id                 int    `json:"summoner1Id"`
	Summoner2Id                 int    `json:"summoner2Id"`
	SummonerId                  string `json:"summonerId"`
	SummonerLevel               int    `json:"summonerLevel"`
	SummonerName                string `json:"summonerName"`
	TeamId                      int    `json:"teamId"`
	TotalDamageDealtToChampions int    `json:"totalDamageDealtToChampions"`
	Win                         bool   `json:"win"`
}

type MatchData struct {
	Metadata struct {
		MatchId string `json:"matchId"`
	} `json:"metadata"`
	Info struct {
		GameCreation       int64         `json:"gameCreation"`
		GameDuration       int           `json:"gameDuration"`
		GameEndTimestamp   int64         `json:"gameEndTimestamp"`
		GameId             int64         `json:"gameId"`
		GameMode           string        `json:"gameMode"`
		GameName           string        `json:"gameName"`
		GameStartTimestamp int64         `json:"gameStartTimestamp"`
		GameType           string        `json:"gameType"`
		GameVersion        string        `json:"gameVersion"`
		MapId              int           `json:"mapId"`
		Participants       []Participant `json:"participants"`
		PlatformId         string        `json:"platformId"`
		QueueId            int           `json:"queueId"`
	} `json:"info"`
}
