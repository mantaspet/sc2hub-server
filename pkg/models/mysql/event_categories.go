package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type EventCategoryModel struct {
	DB *sql.DB
}

func (m *EventCategoryModel) SelectEventCategories() ([]models.EventCategory, error) {
	var ec models.EventCategory
	eventCategories := []models.EventCategory{}
	rows, err := m.DB.Query(`
		SELECT
			id,
		    name,
		    pattern,
		    COALESCE(info_url, '') as info_url,
		    COALESCE(image_url, '') as image_url,
		    COALESCE(description, '') as description,
		    priority
		FROM
		    event_categories
		ORDER BY
		    priority`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ec.ID, &ec.Name, &ec.Pattern, &ec.InfoURL, &ec.ImageURL, &ec.Description, &ec.Priority)
		if err != nil {
			return nil, err
		}
		eventCategories = append(eventCategories, ec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return eventCategories, nil
}

func (m *EventCategoryModel) InsertEventCategory(ec models.EventCategory) (models.EventCategory, error) {
	var res models.EventCategory
	var maxPriority int
	err := m.DB.QueryRow(`SELECT MAX(priority) FROM event_categories`).Scan(&maxPriority)
	qRes, err := m.DB.Exec(`
		INSERT INTO
		  	event_categories (name, pattern, info_url, image_url, description, priority)
		VALUES
		    (?, ?, ?, ?, ?, ?)`, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, &ec.Description, maxPriority+1)
	if err != nil {
		return res, err
	}
	id, err := qRes.LastInsertId()
	if err != nil {
		return res, err
	}
	row := m.DB.QueryRow(`
		SELECT
			id,
		    name,
		    pattern,
		    COALESCE(info_url, '') as info_url,
		    COALESCE(image_url, '') as image_url,
		    COALESCE(description, '') as description,
		    priority
		FROM
		    event_categories
		WHERE
		    id=?`, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.ImageURL, &res.Description, &res.Priority); err != nil {
		return res, err
	}
	return res, nil
}

func (m *EventCategoryModel) UpdateEventCategory(id string, ec models.EventCategory) (models.EventCategory, error) {
	var res models.EventCategory
	_, err := m.DB.Exec(`
		UPDATE
		  	event_categories
		SET
		    name=?,
		    pattern=?,
		    info_url=?,
		    image_url=?,
		    description=?
		WHERE
		    id=?`, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, ec.Description, id)
	if err != nil {
		return res, err
	}
	row := m.DB.QueryRow(`
		SELECT
		    id,
		    name,
		    pattern,
		    COALESCE(info_url, '') as info_url,
		    COALESCE(image_url, '') as image_url,
		    COALESCE(description, '') as description,
		    priority
		FROM
		    event_categories
		WHERE
		    id=?`, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.ImageURL, &res.Description, &res.Priority); err != nil {
		return res, err
	}
	return res, nil
}

func (m *EventCategoryModel) DeleteEventCategory(id string) error {
	var oldPrio int
	var maxPrio int
	tx, err := m.DB.Begin()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.QueryRow(`SELECT priority FROM event_categories WHERE id=?`, id).Scan(&oldPrio)
	_ = tx.QueryRow(`SELECT max(priority) as max FROM event_categories`).Scan(&maxPrio)
	for i := oldPrio + 1; i <= maxPrio; i++ {
		_, err = tx.Exec(`UPDATE event_categories SET priority=? WHERE priority=?`, i-1, i)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	res, err := tx.Exec(`
		DELETE FROM
			event_categories
		WHERE
		    id=?`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	} else if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *EventCategoryModel) UpdateEventCategoryPriorities(id int, newPrio int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var oldPrio int
	_ = tx.QueryRow(`SELECT priority FROM event_categories WHERE id=?`, id).Scan(&oldPrio)
	if oldPrio == newPrio {
		return nil
	}
	if newPrio > oldPrio {
		for i := oldPrio + 1; i <= newPrio; i++ {
			_, err = tx.Exec(`UPDATE event_categories SET priority=? WHERE priority=?`, i-1, i)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	} else {
		for i := oldPrio - 1; i >= newPrio; i-- {
			_, err = tx.Exec(`UPDATE event_categories SET priority=? WHERE priority=?`, i+1, i)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}
	_, err = tx.Exec(`UPDATE event_categories SET priority=? WHERE id=?`, newPrio, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *EventCategoryModel) AssignCategories(events []models.Event) ([]models.Event, error) {
	eventsWithCategories := make([]models.Event, 0, len(events))
	eventCategories, err := m.SelectEventCategories()
	if err != nil {
		return nil, err
	}
	for _, e := range events {
		for _, ec := range eventCategories {
			if strings.Contains(strings.ToLower(e.Title), ec.Pattern) {
				e.EventCategoryID = ec.ID
				break
			}
		}
		eventsWithCategories = append(eventsWithCategories, e)
	}
	return eventsWithCategories, nil
}

func (m *EventCategoryModel) LoadCategories(events []models.Event) ([]models.Event, error) {
	eventCategories, err := m.SelectEventCategories()
	eventsWithCategories := make([]models.Event, 0, len(events))
	if err != nil {
		return nil, err
	}
	for _, e := range events {
		for _, ec := range eventCategories {
			if e.EventCategoryID == ec.ID {
				e.EventCategory = ec
				break
			}
		}
		eventsWithCategories = append(eventsWithCategories, e)
	}
	return eventsWithCategories, nil
}

func (m *EventCategoryModel) fieldExists(table string, field string, value interface{}, id int) error {
	var res interface{}
	query := fmt.Sprintf("SELECT NULL FROM %s WHERE %s='%v'", table, field, value)
	if id > 0 {
		query += fmt.Sprintf(" AND id<>%v", id)
	}
	fmt.Println(query)
	err := m.DB.QueryRow(query).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return nil
	}
	return errors.New(strings.Title(field) + " must be unique")
}
