package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type TwitchChannelModel struct {
	DB *sql.DB
}

func (m *TwitchChannelModel) SelectAll() ([]*models.TwitchChannel, error) {
	stmt := `
		SELECT id, event_category_id, twitch_user_id, name
		FROM twitch_channels`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tcs := []*models.TwitchChannel{}
	for rows.Next() {
		tc := &models.TwitchChannel{}
		err := rows.Scan(&tc.ID, &tc.EventCategoryID, &tc.TwitchUserID, &tc.Name)
		if err != nil {
			return nil, err
		}
		tcs = append(tcs, tc)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tcs, nil
}
