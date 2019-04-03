package mysql

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type VideoModel struct {
	DB *sql.DB
}

func (m *VideoModel) SelectByCategory(categoryID int) ([]models.Video, error) {
	videos := []models.Video{
		models.Video{
			1, 1, categoryID, "Dark vs Classic GSL Code S Season 1 Ro4", "https://www.twitch.tv/videos/405440725",
		},
		models.Video{
			2, 2, categoryID, "Maru vs Dear GSL Code S Season 1 Ro8", "https://www.twitch.tv/videos/403338946",
		},
		models.Video{
			3, 3, categoryID, "Stats vs Serral PvZ - Grand Final - 2018 WCS Global Finals - StarCraft II", "https://www.youtube.com/watch?v=h0UBfmOJYO4",
		},
	}

	return videos, nil
}
