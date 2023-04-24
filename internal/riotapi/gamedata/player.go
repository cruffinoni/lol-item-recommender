package gamedata

type Player struct {
	SummonerId    string `json:"summonerId" db:"id"`
	SummonerName  string `json:"summonerName" db:"name"`
	SummonerLevel int    `json:"summonerLevel" db:"level"`
	LeaguePoints  int    `json:"leaguePoints"`
	Wins          int    `json:"wins"`
	Losses        int    `json:"losses"`
	AccountId     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
}
