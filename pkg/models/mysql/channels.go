package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type ChannelModel struct {
	DB *sql.DB
}

func (m *ChannelModel) SelectAll() ([]*models.Channel, error) {
	stmt := `
		SELECT id, event_category_id, platform_id, login, title, profile_image_url, pattern
		FROM channels`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := []*models.Channel{}
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(&channel.ID, &channel.EventCategoryID, &channel.PlatformID, &channel.Login, &channel.Title, &channel.ProfileImageURL, &channel.Pattern)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (m *ChannelModel) SelectByCategory(categoryID int) ([]*models.Channel, error) {
	stmt := `
		SELECT id, event_category_id, platform_id, login, title, profile_image_url, pattern
		FROM channels
		WHERE event_category_id=?`

	rows, err := m.DB.Query(stmt, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := []*models.Channel{}
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(&channel.ID, &channel.EventCategoryID, &channel.PlatformID, &channel.Login, &channel.Title, &channel.ProfileImageURL, &channel.Pattern)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (m *ChannelModel) Insert(channel models.Channel) (*models.Channel, error) {
	insertStmt := `
		INSERT INTO
		  	channels (id, event_category_id, platform_id, login, title, profile_image_url, pattern)
		VALUES
		    (?, ?, ?, ?, ?, ?, ?)`

	selectStmt := `
		SELECT id, event_category_id, platform_id, login, title, profile_image_url, pattern
		FROM channels
		WHERE id=?`

	_, err := m.DB.Exec(insertStmt, channel.ID, channel.EventCategoryID, channel.PlatformID, channel.Login, channel.Title, channel.ProfileImageURL, channel.Pattern)
	if err != nil {
		return nil, err
	}

	res := &models.Channel{}
	err = m.DB.QueryRow(selectStmt, channel.ID).Scan(&res.ID, &res.EventCategoryID, &res.PlatformID, &res.Login, &res.Title, &res.ProfileImageURL, &res.Pattern)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (m *ChannelModel) Delete(id string) error {
	stmt := `DELETE FROM channels WHERE id=?`

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
