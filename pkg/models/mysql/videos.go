package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type VideoModel struct {
	DB *sql.DB
}

func (m *VideoModel) SelectPage(fromDate string, query string) ([]*models.Video, error) {
	stmt := `SELECT id, COALESCE(event_category_id, 0), platform_id, COALESCE(channel_id, ''), title, duration,
			thumbnail_url, created_at
	  	FROM videos`

	if fromDate == "" {
		stmt += " WHERE 1<>?"
	} else {
		stmt += " WHERE created_at<=?"
	}

	if query == "" {
		stmt += " AND 1<>?"
	} else {
		stmt += " AND title LIKE ?"
	}

	stmt += ` ORDER BY created_at DESC LIMIT ?`

	rows, err := m.DB.Query(stmt, fromDate, "%"+query+"%", models.VideoPageLength+1)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos, err := parseVideoRows(rows)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) SelectRecent() ([]*models.Video, error) {
	stmt := `SELECT id, COALESCE(event_category_id, 0), platform_id, COALESCE(channel_id, ''), title, duration,
			COALESCE(thumbnail_url, ''), created_at
	  	FROM videos
	  	ORDER BY created_at DESC 
	  	LIMIT 16`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos, err := parseVideoRows(rows)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) SelectEventBroadcasts(categoryID int, date string) ([]*models.Video, error) {
	stmt := `SELECT id, COALESCE(event_category_id, 0), platform_id, COALESCE(channel_id, ''), title, duration,
			COALESCE(thumbnail_url, ''), created_at
	  	FROM videos
	  	WHERE event_category_id=? AND created_at LIKE ? AND type='archive'
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, categoryID, "%"+date+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos, err := parseVideoRows(rows)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) SelectByCategory(categoryID int, query string) ([]*models.Video, error) {
	stmt := `SELECT id, COALESCE(event_category_id, 0), platform_id, COALESCE(channel_id, ''), title, duration,
			COALESCE(thumbnail_url, ''), created_at
	  	FROM videos
	  	WHERE event_category_id=? AND title LIKE ?
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, categoryID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos, err := parseVideoRows(rows)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) SelectByPlayer(playerID int, query string) ([]*models.Video, error) {
	stmt := `
		SELECT videos.id, COALESCE(videos.event_category_id, 0), videos.platform_id, COALESCE(videos.channel_id, ''),
			videos.title, videos.duration, COALESCE(videos.thumbnail_url, ''), videos.created_at
		FROM videos
		INNER JOIN player_videos
		ON player_videos.video_id=videos.id
		WHERE player_videos.player_id=? AND title LIKE ?`

	rows, err := m.DB.Query(stmt, playerID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos, err := parseVideoRows(rows)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) InsertOrUpdateMany(videos []*models.Video) (int64, error) {
	valueStrings := make([]string, 0, len(videos))
	valueArgs := make([]interface{}, 0, len(videos)*9)
	for _, v := range videos {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, v.ID)
		valueArgs = append(valueArgs, v.EventCategoryID)
		valueArgs = append(valueArgs, v.PlatformID)
		valueArgs = append(valueArgs, v.ChannelID)
		valueArgs = append(valueArgs, v.Title)
		valueArgs = append(valueArgs, v.Duration)
		valueArgs = append(valueArgs, v.ThumbnailURL)
		valueArgs = append(valueArgs, v.Type)
		valueArgs = append(valueArgs, v.CreatedAt)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO videos(id, event_category_id, platform_id, channel_id, title, duration, thumbnail_url, type,
			created_at)
		VALUES %s 
		ON DUPLICATE KEY UPDATE title=VALUES(title), duration=VALUES(duration), thumbnail_url=VALUES(thumbnail_url),
			created_at=VALUES(created_at);`,
		strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}

	return rowCnt, nil
}

func parseVideoRows(rows *sql.Rows) ([]*models.Video, error) {
	videos := []*models.Video{}
	for rows.Next() {
		v := &models.Video{}
		err := rows.Scan(&v.ID, &v.EventCategoryID, &v.PlatformID, &v.ChannelID, &v.Title, &v.Duration,
			&v.ThumbnailURL, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}
	return videos, nil
}
