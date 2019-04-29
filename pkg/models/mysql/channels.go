package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type ChannelModel struct {
	DB *sql.DB
}

func (m *ChannelModel) SelectAllFromTwitch() ([]*models.Channel, error) {
	stmt := `
		SELECT id, login, platform_id
		FROM channels
		WHERE platform_id=1`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	channels := []*models.Channel{}
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(&channel.ID, &channel.Login, &channel.PlatformID)
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

func (m *ChannelModel) SelectFromAllCategories(platformID int) ([]*models.Channel, error) {
	stmt := `
		SELECT channels.id, channels.platform_id, event_categories.id, event_categories.pattern
		FROM event_category_channels
		INNER JOIN channels
		ON event_category_channels.channel_id = channels.id
		INNER JOIN event_categories
		ON event_category_channels.event_category_id = event_categories.id`

	if platformID > 0 {
		stmt += " WHERE platform_id=?"
	} else {
		stmt += " WHERE -1<>?"
	}

	stmt += ` ORDER BY event_categories.priority`

	rows, err := m.DB.Query(stmt, platformID)
	if err != nil {
		return nil, err
	}

	channels := []*models.Channel{}
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(&channel.ID, &channel.PlatformID, &channel.EventCategoryID, &channel.Pattern)
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

func (m *ChannelModel) SelectByCategory(categoryID int, platformID int) ([]*models.Channel, error) {
	stmt := `
		SELECT channels.id, channels.platform_id, channels.login, channels.title,
		       COALESCE(channels.profile_image_url, '')
		FROM channels
		INNER JOIN event_category_channels
		ON event_category_channels.channel_id=channels.id
		WHERE event_category_channels.event_category_id=?`

	var rows *sql.Rows
	var err error

	// platform ID 0 should query all platforms
	if platformID > 0 {
		stmt += " AND platform_id=?"
		rows, err = m.DB.Query(stmt, categoryID, platformID)
	} else {
		rows, err = m.DB.Query(stmt, categoryID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := []*models.Channel{}
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(&channel.ID, &channel.PlatformID, &channel.Login, &channel.Title, &channel.ProfileImageURL)
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

func (m *ChannelModel) Insert(channel models.Channel, categoryID int) (*models.Channel, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}

	insertStmt := `INSERT INTO	channels (id, platform_id, login, title, profile_image_url)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE title=VALUES(title), profile_image_url=VALUES(profile_image_url);`

	_, err = tx.Exec(insertStmt, channel.ID, channel.PlatformID, channel.Login, channel.Title, channel.ProfileImageURL)
	if err != nil {
		return nil, err
	}

	insertStmt = `INSERT INTO event_category_channels (event_category_id, channel_id) VALUES (?, ?)`
	_, err = tx.Exec(insertStmt, categoryID, channel.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	selectStmt := `
		SELECT id, platform_id, login, title, COALESCE(profile_image_url, '')
		FROM channels
		WHERE id=?`

	res := &models.Channel{}
	err = tx.QueryRow(selectStmt, channel.ID).Scan(&res.ID, &res.PlatformID, &res.Login, &res.Title, &res.ProfileImageURL)
	if err != nil {
		_ = tx.Rollback()
		return res, err
	}

	err = tx.Commit()
	return res, err
}

func (m *ChannelModel) DeleteFromCategory(channelID string, categoryID int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	stmt := `DELETE FROM event_category_channels WHERE channel_id=? AND event_category_id=?`
	res, err := tx.Exec(stmt, channelID, categoryID)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if rowCnt == 0 {
		return models.ErrNotFound
	} else if err != nil {
		return err
	}

	var count int64
	stmt = `SELECT COUNT(*) FROM event_category_channels WHERE channel_id=?`
	err = tx.QueryRow(stmt, channelID).Scan(&count)
	if count > 0 {
		return tx.Commit()
	}

	stmt = `DELETE FROM channels WHERE id=?`
	res, err = tx.Exec(stmt, channelID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	rowCnt, err = res.RowsAffected()
	if rowCnt == 0 {
		_ = tx.Rollback()
		return models.ErrNotFound
	} else if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
