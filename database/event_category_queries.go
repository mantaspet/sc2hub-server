package database

import (
	"database/sql"
)

func SelectEventCategories() ([]EventCategory, error) {
	var ec EventCategory
	eventCategories := []EventCategory{}
	rows, err := db.Query(`
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

func InsertEventCategory(ec EventCategory) (EventCategory, error) {
	var res EventCategory
	var maxPriority int
	err := db.QueryRow(`SELECT MAX(priority) FROM event_categories`).Scan(&maxPriority)
	qRes, err := db.Exec(`
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
	row := db.QueryRow(`
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

func UpdateEventCategory(id string, ec EventCategory) (EventCategory, error) {
	var res EventCategory
	_, err := db.Exec(`
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
	row := db.QueryRow(`
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

func DeleteEventCategory(id string) error {
	var oldPrio int
	var maxPrio int
	tx, err := db.Begin()
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

func UpdateEventCategoryPriorities(id int, newPrio int) error {
	tx, err := db.Begin()
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
