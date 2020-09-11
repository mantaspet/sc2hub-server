package main

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"log"
	"net/http"
)

type application struct {
	httpClient        *http.Client
	db                *sql.DB // TODO find a better solution. This is used only in pkg validators SQLUnique function
	appOrigin         string
	adminOrigin       string
	twitchAccessToken string
	errorLog          *log.Logger
	infoLog           *log.Logger
	twitchGameId      int
	events            interface {
		SelectInDateRange(dateFrom string, dateTo string) ([]*models.Event, error)
		SelectOne(id string) (*models.Event, error)
		InsertMany(events []models.Event) (int64, error)
	}
	eventCategories interface {
		SelectAll() ([]*models.EventCategory, error)
		SelectAllPatterns() ([]*models.EventCategory, error)
		SelectOne(id string) (*models.EventCategory, error)
		Insert(ec models.EventCategory) (*models.EventCategory, error)
		Update(id string, ec models.EventCategory) (*models.EventCategory, error)
		Delete(id string) error
		UpdatePriorities(id int, newPrio int) error
		AssignToEvents(events []models.Event) ([]models.Event, error)
		InsertEventCategoryArticles(ecArticles []models.EventCategoryArticle) (int64, error)
	}
	players interface {
		SelectPage(fromID int, query string) ([]*models.Player, error)
		SelectOne(id int) (*models.Player, error)
		SelectAllPlayerIDs() ([]string, error)
		SelectAllPlayerIDsAndIDs() ([]*models.Player, error)
		InsertMany(players []models.Player) (int64, error)
		InsertPlayerVideos(playerVideos []models.PlayerVideo) (int64, error)
		InsertPlayerArticles(playerArticles []models.PlayerArticle) (int64, error)
	}
	videos interface {
		SelectPage(pageSize int, from int, query string) ([]*models.Video, error)
		SelectRecent() ([]*models.Video, error)
		SelectEventBroadcasts(categoryID int, date string) ([]*models.Video, error)
		SelectByCategory(pageSize int, from int, query string, categoryID int) ([]*models.Video, error)
		SelectByPlayer(pageSize int, from int, query string, playerID int) ([]*models.Video, error)
		InsertOrUpdateMany(videos []*models.Video) (int64, error)
		UpdateMetadata(videos []*models.Video) error
	}
	articles interface {
		SelectPage(pageSize int, from int, query string) ([]*models.Article, error)
		SelectByCategory(pageSize int, from int, query string, categoryID int) ([]*models.Article, error)
		SelectByPlayer(pageSize int, from int, query string, playerID int) ([]*models.Article, error)
		SelectLastInserted(amount int64) ([]*models.Article, error)
		InsertMany(articles []models.Article) (int64, error)
	}
	channels interface {
		SelectFromAllCategories(platformID int) ([]*models.Channel, error)
		SelectAllFromTwitch() ([]*models.Channel, error)
		SelectByCategory(categoryID int, platformID int) ([]*models.Channel, error)
		Insert(channel models.Channel, categoryID int) (*models.Channel, error)
		DeleteFromCategory(channelID string, categoryID int) error
	}
	users interface {
		SelectOne(username string) (*models.User, error)
	}
}
