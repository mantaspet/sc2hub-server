package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strconv"
	"strings"
)

type VideoModel struct {
	DB *sql.DB
}

func (m *VideoModel) SelectByCategory(categoryID string) ([]*models.Video, error) {
	id, _ := strconv.Atoi(categoryID)
	stmt := `SELECT
			id,
			COALESCE(event_id, 0) as event_id,
			COALESCE(event_category_id, 0) as event_category_id,
			COALESCE(channel_id, 0) as channel_id,
			COALESCE(twitch_id, 0) as twitch_id,
			title,
			duration,
			created_at
	  	FROM videos
	  	WHERE event_category_id=?
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []*models.Video{}
	for rows.Next() {
		v := &models.Video{}
		err := rows.Scan(&v.ID, &v.EventID, &v.EventCategoryID, &v.ChannelID, &v.TwitchID, &v.Title, &v.Duration, &v.CreatedAt)
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
	valueArgs := make([]interface{}, 0, len(videos)*6)
	for _, v := range videos {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, v.EventCategoryID)
		valueArgs = append(valueArgs, v.ChannelID)
		valueArgs = append(valueArgs, v.TwitchID)
		valueArgs = append(valueArgs, v.Title)
		valueArgs = append(valueArgs, v.Duration)
		valueArgs = append(valueArgs, v.CreatedAt)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO videos(event_category_id, channel_id, twitch_id, title, duration, created_at)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			title=VALUES(title);`, strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	_, _ = m.DB.Exec(`ALTER TABLE videos AUTO_INCREMENT=1`) // to prevent ON DUPLICATE KEY triggers from inflating next ID
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}

	return rowCnt, nil
}
