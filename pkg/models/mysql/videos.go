package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strconv"
	"strings"
	"time"
)

type VideoModel struct {
	DB *sql.DB
}

func (m *VideoModel) SelectByCategory(categoryID string) ([]models.Video, error) {
	id, _ := strconv.Atoi(categoryID)
	videos := []models.Video{
		models.Video{
			1, 1, id, 1, 1, "Dark vs Classic GSL Code S Season 1 Ro4", "1h20min30s", time.Now(),
		},
		models.Video{
			2, 2, id, 1, 1, "Maru vs Dear GSL Code S Season 1 Ro8", "1h20min30s", time.Now(),
		},
		models.Video{
			3, 3, id, 1, 1, "Stats vs Serral PvZ - Grand Final - 2018 WCS Global Finals - StarCraft II", "1h20min30s", time.Now(),
		},
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
