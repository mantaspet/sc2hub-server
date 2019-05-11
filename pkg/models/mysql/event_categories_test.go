package mysql

import (
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"reflect"
	"strings"
	"testing"
)

func TestEventCategoryModel_SelectAll(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name    string
		wantRes []*models.EventCategory
		wantErr error
	}{
		{
			name:    "Select all",
			wantRes: mock.EventCategories,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			res, err := m.SelectAll()

			if err != tt.wantErr {
				t.Errorf("want %v; got %s", tt.wantErr, err)
			}

			if !reflect.DeepEqual(res, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, res)
			}
		})
	}
}

func TestEventCategoryModel_SelectAllPatterns(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name    string
		wantRes []*models.EventCategory
		wantErr error
	}{
		{
			name:    "Select patterns",
			wantRes: mock.EventCategoryPatterns,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			res, err := m.SelectAllPatterns()

			if err != tt.wantErr {
				t.Errorf("want %v; got %s", tt.wantErr, err)
			}

			if !reflect.DeepEqual(res, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, res)
			}
		})
	}
}

func TestEventCategoryModel_SelectOne(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name    string
		id      string
		wantRes *models.EventCategory
		wantErr error
	}{
		{
			name:    "Valid ID",
			id:      "1",
			wantRes: mock.EventCategories[0],
			wantErr: nil,
		},
		{
			name:    "Zero ID",
			id:      "0",
			wantRes: nil,
			wantErr: models.ErrNotFound,
		},
		{
			name:    "Non-existent ID",
			id:      "5",
			wantRes: nil,
			wantErr: models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			ec, err := m.SelectOne(tt.id)

			if err != tt.wantErr {
				t.Errorf("want %v; got %s", tt.wantErr, err)
			}

			if !reflect.DeepEqual(ec, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, ec)
			}
		})
	}
}

func TestEventCategoryModel_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	newCategory := &models.EventCategory{
		ID: 4, Name: "World Electronic Sports Games", Pattern: "wesg", InfoURL: "https://infourl.com",
		ImageURL: "http://imageurl.com", Description: "", Priority: 4,
	}

	tests := []struct {
		name    string
		ec      *models.EventCategory
		wantRes *models.EventCategory
		wantErr error
	}{
		{
			name:    "Duplicate field",
			ec:      mock.EventCategories[0],
			wantRes: nil,
			wantErr: errors.New("'event_categories_name_uindex'"),
		},
		{
			name:    "Valid",
			ec:      newCategory,
			wantRes: newCategory,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			ec, err := m.Insert(*tt.ec)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(ec, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, ec)
			}
		})
	}
}

func TestEventCategoryModel_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	updatedCategory := &models.EventCategory{
		ID: 3, Name: "World Electronic Sports Games", Pattern: "wesg", InfoURL: "https://infourl.com",
		ImageURL: "http://imageurl.com", Description: "", Priority: 3,
	}

	tests := []struct {
		name    string
		id      string
		ec      *models.EventCategory
		wantRes *models.EventCategory
		wantErr error
	}{
		{
			name:    "Duplicate field",
			ec:      mock.EventCategories[0],
			id:      "3",
			wantRes: nil,
			wantErr: errors.New("'event_categories_name_uindex'"),
		},
		{
			name:    "Valid",
			id:      "3",
			ec:      updatedCategory,
			wantRes: updatedCategory,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			ec, err := m.Update(tt.id, *tt.ec)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(ec, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, ec)
			}
		})
	}
}

func TestEventCategoryModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "With existing events",
			id:      "1",
			wantErr: errors.New("a foreign key constraint fails"),
		},
		{
			name:    "Valid",
			id:      "3",
			wantErr: nil,
		},
		{
			name:    "Not existing",
			id:      "4",
			wantErr: errors.New("no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			err := m.Delete(tt.id)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}
		})
	}
}

func TestEventCategoryModel_UpdatePriorities(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name    string
		id      int
		newPrio int
		wantErr error
	}{
		{
			name:    "Valid - move forward",
			id:      1,
			newPrio: 3,
			wantErr: nil,
		},
		{
			name:    "Valid - move backward",
			id:      3,
			newPrio: 2,
			wantErr: nil,
		},
		{
			name:    "Valid - stay",
			id:      3,
			newPrio: 3,
			wantErr: nil,
		},
		{
			name:    "Non existing category",
			id:      4,
			newPrio: 3,
			wantErr: errors.New("no rows in result set"),
		},
		{
			name:    "Invalid prio",
			id:      3,
			newPrio: -1,
			wantErr: errors.New("of range value for column 'priority'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			err := m.UpdatePriorities(tt.id, tt.newPrio)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}
		})
	}
}

func TestEventCategoryModel_AssignToEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	eventsWithCategories := make([]models.Event, 0, 3)
	eventWithoutCategories := make([]models.Event, 0, 3)
	for _, e := range mock.Events {
		val := *e
		val.EventCategoryID = 0
		eventsWithCategories = append(eventsWithCategories, *e)
		eventWithoutCategories = append(eventWithoutCategories, val)
	}

	tests := []struct {
		name    string
		id      int
		events  []models.Event
		wantRes []models.Event
		wantErr error
	}{
		{
			name:    "With data",
			events:  eventWithoutCategories,
			wantRes: eventsWithCategories,
			wantErr: nil,
		},
		{
			name:    "Without data",
			events:  []models.Event{},
			wantRes: []models.Event{},
			wantErr: nil,
		},
		//{
		//	name: "Invalid prio",
		//	id: 3,
		//	newPrio: -1,
		//	wantErr: errors.New("of range value for column 'priority'"),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			res, err := m.AssignToEvents(tt.events)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(res, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, res)
			}
		})
	}
}

func TestEventCategoryModel_InsertEventCategoryArticles(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name       string
		id         int
		ecArticles []models.EventCategoryArticle
		wantRes    int64
		wantErr    error
	}{
		{
			name:       "With data",
			ecArticles: mock.EventCategoryArticles,
			wantRes:    2,
			wantErr:    nil,
		},
		{
			name:       "Without data",
			ecArticles: []models.EventCategoryArticle{},
			wantRes:    0,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			res, err := m.InsertEventCategoryArticles(tt.ecArticles)

			if err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) == false {
				t.Errorf("want %v; in %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(res, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, res)
			}
		})
	}
}
