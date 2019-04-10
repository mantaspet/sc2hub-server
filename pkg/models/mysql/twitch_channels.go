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
		SELECT id, event_category_id, twitch_user_id, login, display_name, profile_image_url
		FROM twitch_channels`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tcs := []*models.TwitchChannel{}
	for rows.Next() {
		tc := &models.TwitchChannel{}
		err := rows.Scan(&tc.ID, &tc.EventCategoryID, &tc.TwitchUserID, &tc.Login, &tc.DisplayName, &tc.ProfileImageURL)
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

func (m *TwitchChannelModel) SelectByCategory(categoryID int) ([]*models.TwitchChannel, error) {
	stmt := `
		SELECT id, event_category_id, twitch_user_id, login, display_name, profile_image_url
		FROM twitch_channels
		WHERE event_category_id=?`

	rows, err := m.DB.Query(stmt, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tcs := []*models.TwitchChannel{}
	for rows.Next() {
		tc := &models.TwitchChannel{}
		err := rows.Scan(&tc.ID, &tc.EventCategoryID, &tc.TwitchUserID, &tc.Login, &tc.DisplayName, &tc.ProfileImageURL)
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

func (m *TwitchChannelModel) Insert(tc models.TwitchChannel) (*models.TwitchChannel, error) {
	insertStmt := `
		INSERT INTO
		  	twitch_channels (event_category_id, twitch_user_id, login, display_name, profile_image_url)
		VALUES
		    (?, ?, ?, ?, ?)`

	selectStmt := `
		SELECT id, event_category_id, twitch_user_id, login, display_name, profile_image_url
		FROM twitch_channels
		WHERE id=?`

	insertRes, err := m.DB.Exec(insertStmt, tc.EventCategoryID, tc.TwitchUserID, tc.Login, tc.DisplayName, tc.ProfileImageURL)
	if err != nil {
		return nil, err
	}
	id, err := insertRes.LastInsertId()
	if err != nil {
		return nil, err
	}

	res := &models.TwitchChannel{}
	err = m.DB.QueryRow(selectStmt, id).Scan(&res.ID, &res.EventCategoryID, &res.TwitchUserID, &res.Login, &res.DisplayName, &res.ProfileImageURL)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (m *TwitchChannelModel) Delete(id int) error {
	stmt := `DELETE FROM twitch_channels WHERE id=?`

	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if rowCnt == 0 {
		return models.ErrNotFound
	} else if err != nil {
		return err
	}

	return err
}
