package main

import (
	"database/sql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"log"
)

type application struct {
	db       *sql.DB // TODO find a better solution. This is used only in pkg validators SQLUnique function
	origin   string
	errorLog *log.Logger
	infoLog  *log.Logger
	events   interface {
		SelectInDateRange(dateFrom string, dateTo string) ([]*models.Event, error)
		SelectOne(id string) (*models.Event, error)
		InsertMany(events []models.Event) (int64, error)
	}
	eventCategories interface {
		SelectAll() ([]*models.EventCategory, error)
		SelectOne(id string) (*models.EventCategory, error)
		Insert(ec models.EventCategory) (*models.EventCategory, error)
		Update(id string, ec models.EventCategory) (*models.EventCategory, error)
		Delete(id string) error
		UpdatePriorities(id int, newPrio int) error
		AssignToEvents(events []models.Event) ([]models.Event, error)
		LoadOnEvents(events []*models.Event) ([]*models.Event, error)
	}
	players interface {
		SelectAllPlayers() ([]*models.Player, error)
		InsertMany(players []models.Player) (int64, error)
	}
	videos interface {
		SelectByCategory(categoryID string) ([]models.Video, error)
		InsertOrUpdateMany(videos []*models.Video) (int64, error)
	}
	articles interface {
		SelectByCategory(categoryID string) ([]models.Article, error)
	}
	twitchChannels interface {
		SelectAll() ([]*models.TwitchChannel, error)
	}
}
