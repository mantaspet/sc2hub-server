package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type PlayerModel struct {
	DB *sql.DB
}

func (m *PlayerModel) SelectAllPlayers() ([]*models.Player, error) {
	stmt := `SELECT
			id, player_id, name, race, team, country, total_earnings,
       		COALESCE(date_of_birth, '') as date_of_birth, liquipedia_url, image_url, stream_url, is_retired
	  	FROM players`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []*models.Player{}
	for rows.Next() {
		p := &models.Player{}
		err := rows.Scan(&p.ID, &p.PlayerID, &p.Name, &p.Race, &p.Team, &p.Country, &p.TotalEarnings, &p.DateOfBirth, &p.LiquipediaURL, &p.ImageURL, &p.StreamURL, &p.IsRetired)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (m *PlayerModel) InsertMany(players []models.Player) (int64, error) {
	valueStrings := make([]string, 0, len(players))
	valueArgs := make([]interface{}, 0, len(players)*10)
	for _, p := range players {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, p.PlayerID)
		valueArgs = append(valueArgs, p.Name)
		fmt.Println(p.Name)
		valueArgs = append(valueArgs, p.Race)
		valueArgs = append(valueArgs, p.Team)
		valueArgs = append(valueArgs, p.Country)
		valueArgs = append(valueArgs, p.LiquipediaURL)
		valueArgs = append(valueArgs, p.StreamURL)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO players(player_id, name, race, team, country, liquipedia_url, stream_url)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			name=VALUES(name),
			race=VALUES(race),
			team=VALUES(team),
			country=VALUES(country),
			liquipedia_url=VALUES(liquipedia_url),
			stream_url=VALUES(stream_url);`, strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	_, _ = m.DB.Exec(`ALTER TABLE players AUTO_INCREMENT=1`) // to prevent ON DUPLICATE KEY triggers from inflating next ID
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}

	return rowCnt, nil
}
