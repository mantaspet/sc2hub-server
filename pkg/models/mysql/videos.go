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

func (m *VideoModel) SelectByCategory(categoryID int, query string) ([]*models.Video, error) {
	stmt := `SELECT
			id,
			COALESCE(event_id, 0) as event_id,
			COALESCE(event_category_id, 0) as event_category_id,
			platform_id,
			COALESCE(channel_id, '') as channel_id,
			title,
			duration,
			created_at
	  	FROM videos
	  	WHERE event_category_id=? AND title LIKE ?
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, categoryID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []*models.Video{}
	for rows.Next() {
		v := &models.Video{}
		err := rows.Scan(&v.ID, &v.EventID, &v.EventCategoryID, &v.PlatformID, &v.ChannelID, &v.Title, &v.Duration, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (m *VideoModel) InsertOrUpdateMany(videos []*models.Video) (int64, error) {
	valueStrings := make([]string, 0, len(videos))
	valueArgs := make([]interface{}, 0, len(videos)*7)
	for _, v := range videos {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, v.ID)
		valueArgs = append(valueArgs, v.EventCategoryID)
		valueArgs = append(valueArgs, v.PlatformID)
		valueArgs = append(valueArgs, v.ChannelID)
		valueArgs = append(valueArgs, v.Title)
		valueArgs = append(valueArgs, v.Duration)
		valueArgs = append(valueArgs, v.CreatedAt)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO videos(id, event_category_id, platform_id, channel_id, title, duration, created_at)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			title=VALUES(title);`, strings.Join(valueStrings, ","))

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
