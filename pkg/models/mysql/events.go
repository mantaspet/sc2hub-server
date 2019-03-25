package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) SelectInDateRange(dateFrom string, dateTo string) ([]*models.Event, error) {
	stmt := `SELECT
			id,
			COALESCE(event_category_id, 0) as event_category_id,
			COALESCE(team_liquid_id, 0) as team_liquid_id,
			title,
			stage,
			starts_at,
			info
	  	FROM events
	  	WHERE starts_at BETWEEN ? AND ?
		ORDER BY starts_at`

	rows, err := m.DB.Query(stmt, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*models.Event{}
	for rows.Next() {
		e := &models.Event{}
		err := rows.Scan(&e.ID, &e.EventCategoryID, &e.TeamLiquidID, &e.Title, &e.Stage, &e.StartsAt, &e.Info)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) InsertMany(events []models.Event) (int64, error) {
	valueStrings := make([]string, 0, len(events))
	valueArgs := make([]interface{}, 0, len(events)*5)
	for _, e := range events {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, e.Title)
		if e.EventCategoryID > 0 {
			valueArgs = append(valueArgs, e.EventCategoryID)
		} else {
			valueArgs = append(valueArgs, nil)
		}
		valueArgs = append(valueArgs, e.TeamLiquidID)
		valueArgs = append(valueArgs, e.Stage)
		valueArgs = append(valueArgs, e.StartsAt)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO events(title, event_category_id, team_liquid_id, stage, starts_at)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			title=VALUES(title),
			stage=VALUES(stage),
			starts_at=VALUES(starts_at);`, strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	_, _ = m.DB.Exec(`ALTER TABLE events AUTO_INCREMENT=1`) // to prevent ON DUPLICATE KEY triggers from inflating next ID
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}

	return rowCnt, nil
}
