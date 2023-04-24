package database

import (
	"errors"
	"fmt"
	"strings"

	"LoLItemRecommender/internal/riotapi/gamedata"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const DSN = "root:password@tcp(127.0.0.1:3306)/lol-item-recommender"

type DB struct {
	db *sqlx.DB
}

func (d *DB) createTableMatches() error {
	query := `
		CREATE TABLE IF NOT EXISTS matches (
		    id BIGINT PRIMARY KEY,
		    match_uid VARCHAR(255) NOT NULL,
		    creation BIGINT NOT NULL,
		    duration INT NOT NULL,
		    end_timestamp BIGINT NOT NULL,
		    mode VARCHAR(255) NOT NULL,
		    name VARCHAR(255) NOT NULL,
		    start_timestamp BIGINT NOT NULL,
		    type VARCHAR(255) NOT NULL,
		    version VARCHAR(255) NOT NULL,
		    map_id INT NOT NULL,
		    platform_id VARCHAR(255) NOT NULL,
		    queue_id INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) createTableParticipants() error {
	query := `
		CREATE TABLE IF NOT EXISTS participants (
		    participant_id INT NOT NULL,
		    match_id BIGINT NOT NULL,
		    summoner_id VARCHAR(255) NOT NULL,
		    champion_id INT NOT NULL,
		    team_id INT NOT NULL,
		    role VARCHAR(255) NOT NULL,
		    lane VARCHAR(255) NOT NULL,
		    kills INT NOT NULL,
		    deaths INT NOT NULL,
		    assists INT NOT NULL,
		    champ_level INT NOT NULL,
		    total_damage_dealt_to_champions INT NOT NULL,
		    gold_earned INT NOT NULL,
		    win BOOLEAN NOT NULL,
		    summoner_spell1_id INT NOT NULL,
		    summoner_spell2_id INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		    PRIMARY KEY (participant_id, match_id),
		    FOREIGN KEY (match_id) REFERENCES matches(id),
		    FOREIGN KEY (summoner_id) REFERENCES summoners(id)
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) createTableItems() error {
	query := `
		CREATE TABLE IF NOT EXISTS items (
		    match_id BIGINT NOT NULL,
		    participant_id INT NOT NULL,
		    item_id INT NOT NULL,
		    slot INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		    PRIMARY KEY (participant_id, match_id, item_id),
		    FOREIGN KEY (match_id) REFERENCES matches(id)
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) createTableSummoners() error {
	query := `
		CREATE TABLE IF NOT EXISTS summoners (
		    id VARCHAR(255) PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    level INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) createTablePerks() error {
	query := `
		CREATE TABLE IF NOT EXISTS perks (
		    match_id BIGINT NOT NULL,
		    participant_id INT NOT NULL,
		    style INT NOT NULL,
		    perk INT NOT NULL,
		    var1 INT NOT NULL,
		    var2 INT NOT NULL,
		    var3 INT NOT NULL,
		    PRIMARY KEY (participant_id, match_id),
		    FOREIGN KEY (match_id) REFERENCES matches(id)
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) createTableStatPerks() error {
	query := `
		CREATE TABLE IF NOT EXISTS stat_perks (
		    match_id BIGINT NOT NULL,
		    participant_id INT NOT NULL,
		    defense INT NOT NULL,
		    flex INT NOT NULL,
		    offense INT NOT NULL,
		    PRIMARY KEY (participant_id, match_id),
		    FOREIGN KEY (match_id) REFERENCES matches(id)
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateTables() error {
	fncs := []func() error{
		d.createTableMatches,
		d.createTableSummoners,
		d.createTableParticipants,
		d.createTableItems,
		d.createTablePerks,
		d.createTableStatPerks,
	}
	for _, f := range fncs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) saveItems(participant *gamedata.Participant, gameID int64) error {
	// Save items information
	for i := 0; i < 7; i++ {
		item := 0
		switch i {
		case 0:
			item = participant.Item0
		case 1:
			item = participant.Item1
		case 2:
			item = participant.Item2
		case 3:
			item = participant.Item3
		case 4:
			item = participant.Item4
		case 5:
			item = participant.Item5
		case 6:
			item = participant.Item6
		}
		if item == 0 {
			continue
		}

		_, err := d.db.Exec(`
				INSERT INTO items (match_id, participant_id, slot, item_id)
				VALUES (?, ?, ?, ?)`,
			gameID, participant.ParticipantId, i, item)

		if err != nil {
			return fmt.Errorf("can't save item: %w", err)
		}
	}
	return nil
}

func (d *DB) savePerks(participant *gamedata.Participant, gameID int64) error {
	// Save perks information
	for _, style := range participant.Perks.Styles {
		for _, selection := range style.Selections {
			_, err := d.db.Exec(`
					INSERT INTO perks (match_id, participant_id, style, perk, var1, var2, var3)
					VALUES (?, ?, ?, ?, ?, ?, ?)
					ON DUPLICATE KEY UPDATE
						perk=VALUES(perk),
						var1=VALUES(var1),
						var2=VALUES(var2),
						var3=VALUES(var3)`,
				gameID, participant.ParticipantId, style.Style, selection.Perk, selection.Var1, selection.Var2, selection.Var3)

			if err != nil {
				return errors.Join(errors.New("can't insert participant style"), err)
			}
		}
	}

	// Save stat perks information
	statPerks := participant.Perks.StatPerks
	_, err := d.db.Exec(`
			INSERT INTO stat_perks (match_id, participant_id, offense, defense, flex)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				offense=VALUES(offense),
				defense=VALUES(defense),
				flex=VALUES(flex)`,
		gameID, participant.ParticipantId, statPerks.Offense, statPerks.Defense, statPerks.Flex)
	return err
}

func (d *DB) saveParticipants(match *gamedata.MatchData) error {
	for _, participant := range match.Info.Participants {
		// Save summoners information
		_, err := d.db.Exec(`
			INSERT INTO summoners (id, name, level)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE
				name=VALUES(name),
				level=VALUES(level)`,
			participant.SummonerId, participant.SummonerName, participant.SummonerLevel)
		if err != nil {
			return fmt.Errorf("can't insert summoner: %w", err)
		}
		if err != nil {
			return fmt.Errorf("can't get last insert id: %w", err)
		}

		_, err = d.db.Exec(`
			INSERT INTO participants (match_id, participant_id, summoner_id, champion_id, team_id, role, lane, kills, deaths, assists, champ_level, total_damage_dealt_to_champions, gold_earned, win, summoner_spell1_id, summoner_spell2_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			match.Info.GameId,
			participant.ParticipantId,
			participant.SummonerId,
			participant.ChampionId,
			participant.TeamId,
			participant.Role,
			participant.Lane,
			participant.Kills,
			participant.Deaths,
			participant.Assists,
			participant.ChampLevel,
			participant.TotalDamageDealtToChampions,
			participant.GoldEarned,
			participant.Win,
			participant.Summoner1Id,
			participant.Summoner2Id)

		if err != nil {
			return fmt.Errorf("can't insert participant: %w", err)
		}

		if err := d.saveItems(&participant, match.Info.GameId); err != nil {
			return fmt.Errorf("can't save items: %w", err)
		}

		if err := d.savePerks(&participant, match.Info.GameId); err != nil {
			return fmt.Errorf("can't insert perks: %w", err)
		}
	}
	return nil
}

func (d *DB) SaveMatch(match *gamedata.MatchData) error {
	// Save match information
	_, err := d.db.Exec(`
		INSERT IGNORE INTO matches (id, match_uid, creation, duration, end_timestamp, mode, name, start_timestamp, type, version, map_id, platform_id, queue_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		match.Info.GameId, match.Metadata.MatchId, match.Info.GameCreation, match.Info.GameDuration, match.Info.GameEndTimestamp, match.Info.GameMode, match.Info.GameName, match.Info.GameStartTimestamp, match.Info.GameType, match.Info.GameVersion, match.Info.MapId, match.Info.PlatformId, match.Info.QueueId)

	if err != nil {
		return fmt.Errorf("can't insert match: %w", err)
	}
	return d.saveParticipants(match)
}

func (d *DB) GetLatestCrawledPlayers() ([]*gamedata.Player, error) {
	var p []*gamedata.Player
	err := d.db.Select(&p, `select id, name, level from summoners order by created_at limit 200`)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *DB) GetAllPlayerNameCrawledExcept(exceptions []string) ([]string, error) {
	var pName []string
	query, args, err := sqlx.In(`SELECT name FROM summoners WHERE name NOT IN (?)`, exceptions)
	if err != nil {
		return nil, err
	}

	query = d.db.Rebind(query)
	err = d.db.Select(&pName, query, args...)

	if err != nil {
		return nil, err
	}
	return pName, nil
}

//func (d *DB) AssociateItemToParticipants(participants []*gamedata.Participant) {
//	items := make(map[int]map[int]*gamedata.Participant)
//	matchIds := make([]string, 0)
//	for participants
//}

func (d *DB) GetMatchesWithChampions(blueTeam, redTeam []*gamedata.ChampionStats) (map[int][]*gamedata.Participant, error) {
	// Prepare the SQL query
	allChamps := make([]string, 0)
	lenTeam1 := len(blueTeam)
	lenTeam2 := len(redTeam)
	sumTeam1ChampsA := ""
	sumTeam1ChampsB := ""
	for i, c := range blueTeam {
		if i > 0 {
			sumTeam1ChampsA += " + "
			sumTeam1ChampsB += " + "
		}
		allChamps = append(allChamps, c.Key)
		sumTeam1ChampsA += fmt.Sprintf("SUM(champion_id = %v AND participant_id BETWEEN 1 AND 5)", c.Key)
		sumTeam1ChampsB += fmt.Sprintf("SUM(champion_id = %v AND participant_id BETWEEN 6 AND 10)", c.Key)
	}

	sumTeam2ChampsA := ""
	sumTeam2ChampsB := ""
	for i, c := range redTeam {
		if i > 0 {
			sumTeam2ChampsA += " + "
			sumTeam2ChampsB += " + "
		}
		allChamps = append(allChamps, c.Key)
		sumTeam2ChampsA += fmt.Sprintf("SUM(champion_id = %v AND participant_id BETWEEN 1 AND 5)", c.Key)
		sumTeam2ChampsB += fmt.Sprintf("SUM(champion_id = %v AND participant_id BETWEEN 6 AND 10)", c.Key)
	}

	query := `
		SELECT p.participant_id,
			   p.match_id,
			   p.summoner_id,
			   p.champion_id,
			   p.team_id,
			   p.role,
			   p.lane,
			   p.kills,
			   p.deaths,
			   p.assists,
			   p.champ_level,
			   p.total_damage_dealt_to_champions,
			   p.gold_earned,
			   p.win,
			   p.summoner_spell1_id,
			   p.summoner_spell2_id,
			   GROUP_CONCAT(DISTINCT i.item_id ORDER BY i.slot SEPARATOR ',') AS items,
			   pe.style,
			   pe.perk,
			   pe.var1,
			   pe.var2,
			   pe.var3,
			   sp.defense,
			   sp.flex,
			   sp.offense
		FROM participants p
				LEFT JOIN items i ON i.match_id = p.match_id AND i.participant_id = p.participant_id
				LEFT JOIN perks pe ON p.match_id = pe.match_id AND p.participant_id = pe.participant_id
				LEFT JOIN stat_perks sp ON p.match_id = sp.match_id AND p.participant_id = sp.participant_id
		WHERE p.match_id IN (
			SELECT match_id
			FROM (
				SELECT match_id,
					   %s AS team1_champs_A,
					   %s AS team2_champs_A,
					   %s AS team1_champs_B,
					   %s AS team2_champs_B
				FROM participants
				WHERE champion_id IN (%s)
				GROUP BY match_id
			) as subquery
			WHERE (team1_champs_A = %d AND team2_champs_B = %d) OR (team2_champs_A = %d AND team1_champs_B = %d)
		);
    `
	query = fmt.Sprintf(query, sumTeam1ChampsA, sumTeam1ChampsB, sumTeam2ChampsA, sumTeam2ChampsB, strings.Join(allChamps, ","), lenTeam1, lenTeam2, lenTeam1, lenTeam2)

	var participants []*Participant
	// Execute the SQL query
	err := d.db.Select(&participants, query)
	if err != nil {
		return nil, err
	}
	participantsPerMatch := make(map[int][]*gamedata.Participant)
	return participantsPerMatch, nil
}

func NewDB() (*DB, error) {
	db, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}
