package mysql

import (
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"reflect"
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

			eventCategory, err := m.SelectOne(tt.id)

			if err != tt.wantErr {
				t.Errorf("want %v; got %s", tt.wantErr, err)
			}

			if !reflect.DeepEqual(eventCategory, tt.wantRes) {
				t.Errorf("want %v; got %v", tt.wantRes, eventCategory)
			}
		})
	}
}
