package database

type Participant struct {
	ID                          int64  `json:"id"`
	ParticipantID               int    `json:"participant_id" db:"participant_id"`
	MatchID                     int64  `json:"match_id" db:"match_id"`
	SummonerID                  string `json:"summoner_id" db:"summoner_id"`
	ChampionID                  int    `json:"champion_id" db:"champion_id"`
	TeamID                      int    `json:"team_id" db:"team_id"`
	Role                        string `json:"role" db:"role"`
	Lane                        string `json:"lane" db:"lane"`
	Kills                       int    `json:"kills" db:"kills"`
	Deaths                      int    `json:"deaths" db:"deaths"`
	Assists                     int    `json:"assists" db:"assists"`
	ChampLevel                  int    `json:"champ_level" db:"champ_level"`
	TotalDamageDealtToChampions int    `json:"total_damage_dealt_to_champions" db:"total_damage_dealt_to_champions"`
	GoldEarned                  int    `json:"gold_earned" db:"gold_earned"`
	Win                         bool   `json:"win" db:"win"`
	SummonerLevel               int    `json:"summoner_level" db:"summoner_level"`
	SummonerSpell1ID            int    `json:"summoner_spell1_id" db:"summoner_spell1_id"`
	SummonerSpell2ID            int    `json:"summoner_spell2_id" db:"summoner_spell2_id"`
}
