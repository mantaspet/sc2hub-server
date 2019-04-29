package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type EventCategoryModel struct {
	DB *sql.DB
}

func (m *EventCategoryModel) SelectAll() ([]*models.EventCategory, error) {
	stmt := `
		SELECT id, name, pattern, info_url, image_url, description, priority
		FROM event_categories
		ORDER BY priority`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventCategories := []*models.EventCategory{}
	for rows.Next() {
		ec := &models.EventCategory{}
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

func (m *EventCategoryModel) SelectAllPatterns() ([]*models.EventCategory, error) {
	stmt := `
		SELECT id, pattern
		FROM event_categories
		ORDER BY priority`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventCategories := []*models.EventCategory{}
	for rows.Next() {
		ec := &models.EventCategory{}
		err := rows.Scan(&ec.ID, &ec.Pattern)
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

func (m *EventCategoryModel) SelectOne(id string) (*models.EventCategory, error) {
	stmt := `
		SELECT id, name, pattern, info_url, image_url, description, priority
		FROM event_categories
		WHERE id=?`

	ec := &models.EventCategory{}
	err := m.DB.QueryRow(stmt, id).Scan(&ec.ID, &ec.Name, &ec.Pattern, &ec.InfoURL, &ec.ImageURL, &ec.Description, &ec.Priority)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return ec, nil
}

func (m *EventCategoryModel) Insert(ec models.EventCategory) (*models.EventCategory, error) {
	insertStmt := `
		INSERT INTO
		  	event_categories (name, pattern, info_url, image_url, description, priority)
		VALUES
		    (?, ?, ?, ?, ?, ?)`

	selectStmt := `
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
		    id=?`

	maxPrioStmt := `SELECT COALESCE(MAX(priority), 1) FROM event_categories`

	var maxPriority int
	err := m.DB.QueryRow(maxPrioStmt).Scan(&maxPriority)
	if err != nil {
		return nil, err
	}

	insertRes, err := m.DB.Exec(insertStmt, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, &ec.Description, maxPriority+1)
	if err != nil {
		return nil, err
	}
	id, err := insertRes.LastInsertId()
	if err != nil {
		return nil, err
	}

	res := &models.EventCategory{}
	err = m.DB.QueryRow(selectStmt, id).Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.ImageURL, &res.Description, &res.Priority)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (m *EventCategoryModel) Update(id string, ec models.EventCategory) (*models.EventCategory, error) {
	updateStmt := `
		UPDATE
		  	event_categories
		SET
		    name=?,
		    pattern=?,
		    info_url=?,
		    image_url=?,
		    description=?
		WHERE
		    id=?`

	selectStmt := `
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
		    id=?`

	updateRes, err := m.DB.Exec(updateStmt, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, ec.Description, id)
	if err != nil {
		return nil, err
	}
	_, err = updateRes.RowsAffected()
	if err != nil {
		return nil, err
	}

	res := &models.EventCategory{}
	err = m.DB.QueryRow(selectStmt, id).Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.ImageURL, &res.Description, &res.Priority)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *EventCategoryModel) Delete(id string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	var currentPrio int
	err = tx.QueryRow(`SELECT priority FROM event_categories WHERE id=?`, id).Scan(&currentPrio)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	var maxPrio int
	err = tx.QueryRow(`SELECT max(priority) as max FROM event_categories`).Scan(&maxPrio)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	updateStmt, err := tx.Prepare(`UPDATE event_categories SET priority=? WHERE priority=?`)
	if err != nil {
		return err
	}
	for i := currentPrio + 1; i <= maxPrio; i++ {
		_, err = updateStmt.Exec(i-1, i)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	res, err := tx.Exec(`DELETE FROM event_categories WHERE id=?`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	rowCnt, err := res.RowsAffected()
	if rowCnt == 0 {
		_ = tx.Rollback()
		return models.ErrNotFound
	} else if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (m *EventCategoryModel) UpdatePriorities(id int, newPrio int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	var currentPrio int
	err = tx.QueryRow(`SELECT priority FROM event_categories WHERE id=?`, id).Scan(&currentPrio)
	if err != nil {
		return nil
	}

	if currentPrio == newPrio {
		return nil
	}

	updateStmt, err := tx.Prepare(`UPDATE event_categories SET priority=? WHERE priority=?`)
	if err != nil {
		return err
	}

	if newPrio > currentPrio {
		for i := currentPrio + 1; i <= newPrio; i++ {
			_, err = updateStmt.Exec(i-1, i)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	} else {
		for i := currentPrio - 1; i >= newPrio; i-- {
			_, err = updateStmt.Exec(i+1, i)
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
	return err
}

func (m *EventCategoryModel) AssignToEvents(events []models.Event) ([]models.Event, error) {
	eventsWithCategories := make([]models.Event, 0, len(events))
	eventCategories, err := m.SelectAll()
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

func (m *EventCategoryModel) InsertEventCategoryArticles(ecArticles []models.EventCategoryArticle) (int64, error) {
	valueStrings := make([]string, 0, len(ecArticles))
	valueArgs := make([]interface{}, 0, len(ecArticles)*2)
	for _, eca := range ecArticles {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, eca.EventCategoryID)
		valueArgs = append(valueArgs, eca.ArticleID)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO event_category_articles(event_category_id, article_id)
		VALUES %s
		ON DUPLICATE KEY UPDATE event_category_id=VALUES(event_category_id)`, strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return rowCnt, err
	}

	return rowCnt, nil
}
