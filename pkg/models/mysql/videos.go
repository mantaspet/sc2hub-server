package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strconv"
	"time"
)

type VideoModel struct {
	DB *sql.DB
}

func (m *VideoModel) SelectByCategory(categoryID string) ([]models.Video, error) {
	id, _ := strconv.Atoi(categoryID)
	videos := []models.Video{
		models.Video{
			1, 1, id, 1, "Dark vs Classic GSL Code S Season 1 Ro4", "https://www.twitch.tv/videos/405440725", "1h20min30s", time.Now(),
		},
		models.Video{
			2, 2, id, 1, "Maru vs Dear GSL Code S Season 1 Ro8", "https://www.twitch.tv/videos/403338946", "1h20min30s", time.Now(),
		},
		models.Video{
			3, 3, id, 1, "Stats vs Serral PvZ - Grand Final - 2018 WCS Global Finals - StarCraft II", "https://www.youtube.com/watch?v=h0UBfmOJYO4", "1h20min30s", time.Now(),
		},
	}

	return videos, nil
}
