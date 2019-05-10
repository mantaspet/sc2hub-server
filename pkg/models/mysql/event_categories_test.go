package mysql

import (
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"reflect"
	"testing"
)

func TestEventCategoryModelGet(t *testing.T) {
	// Skip the test if the `-short` flag is provided when running the test.
	// We'll talk more about this in a moment.
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name              string
		eventCategoryID   string
		wantEventCategory *models.EventCategory
		wantError         error
	}{
		{
			name:            "Valid ID",
			eventCategoryID: "1",
			wantEventCategory: &models.EventCategory{
				ID:          1,
				Name:        "World Championship Series",
				Pattern:     "wcs",
				InfoURL:     "https://liquipedia.net/starcraft2/World_Championship_Series",
				ImageURL:    "https://static-wcs.starcraft2.com/media/images/logo/logo-event-circuit.png",
				Description: "",
				Priority:    4,
			},
			wantError: nil,
		},
		{
			name:              "Zero ID",
			eventCategoryID:   "0",
			wantEventCategory: nil,
			wantError:         models.ErrNotFound,
		},
		{
			name:              "Non-existent ID",
			eventCategoryID:   "20",
			wantEventCategory: nil,
			wantError:         models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := EventCategoryModel{db}

			eventCategory, err := m.SelectOne(tt.eventCategoryID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(eventCategory, tt.wantEventCategory) {
				t.Errorf("want %v; got %v", tt.wantEventCategory, eventCategory)
			}
		})
	}
}
