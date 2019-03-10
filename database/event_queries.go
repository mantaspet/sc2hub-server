package database

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/models"
	"strings"
)

func SelectEvents(dateFrom string, dateTo string) ([]models.Event, error) {
	var event models.Event
	events := []models.Event{}
	rows, err := db.Query(`
		SELECT
			id,
			COALESCE(event_category_id, 0) as event_category_id,
			COALESCE(team_liquid_id, 0) as team_liquid_id,
			title,
			stage,
			starts_at,
			info
	  	FROM events
	  	WHERE starts_at BETWEEN ? AND ?
		ORDER BY starts_at`, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&event.ID, &event.EventCategoryID, &event.TeamLiquidID, &event.Title, &event.Stage, &event.StartsAt, &event.Info)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func InsertEvents(events []models.Event) (int64, error) {
	valueStrings := make([]string, 0, len(events))
	valueArgs := make([]interface{}, 0, len(events)*4)
	for _, e := range events {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, e.Title)
		valueArgs = append(valueArgs, e.TeamLiquidID)
		valueArgs = append(valueArgs, e.Stage)
		valueArgs = append(valueArgs, e.StartsAt)
	}
	q := fmt.Sprintf(`
		INSERT INTO events(title, team_liquid_id, stage, starts_at)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			title=VALUES(title),
			stage=VALUES(stage),
			starts_at=VALUES(starts_at);`, strings.Join(valueStrings, ","))

	res, err := db.Exec(q, valueArgs...)
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}
	_, _ = db.Exec(`ALTER TABLE events AUTO_INCREMENT=1`) // to prevent ON DUPLICATE KEY triggers from inflating next ID
	return rowCnt, nil
}
