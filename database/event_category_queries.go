package database

import (
	"database/sql"
)

func SelectEventCategories() ([]EventCategory, error) {
	var ec EventCategory
	var eventCategories []EventCategory
	rows, err := db.Query(`
		SELECT
			id,
		    name,
		    pattern,
		    COALESCE(info_url, '') as info_url,
		    COALESCE(image_url, '') as image_url,
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
		err := rows.Scan(&ec.ID, &ec.Name, &ec.Pattern, &ec.InfoURL, &ec.ImageURL, &ec.Priority)
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
	qRes, err := db.Exec(`
		INSERT INTO
		  	event_categories (name, pattern, info_url, image_url, priority)
		VALUES
		    (?, ?, ?, ?, ?)`, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, ec.Priority)
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
		    priority
		FROM
		    event_categories
		WHERE
		    id=?`, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.ImageURL, &res.Priority); err != nil {
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
		    priority=?
		WHERE
		    id=?`, ec.Name, ec.Pattern, ec.InfoURL, ec.ImageURL, ec.Priority, id)
	if err != nil {
		return res, err
	}
	row := db.QueryRow(`
		SELECT
		    id,
		    name,
		    pattern,
		    COALESCE(info_url, '') as info_url,
		    priority
		FROM
		    event_categories
		WHERE
		    id=?`, id)
	if err = row.Scan(&res.ID, &res.Name, &res.Pattern, &res.InfoURL, &res.Priority); err != nil {
		return res, err
	}
	return res, nil
}

func DeleteEventCategory(id string) error {
	res, err := db.Exec(`
		DELETE FROM
			event_categories
		WHERE
		    id=?`, id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt == 0 {
		return sql.ErrNoRows
	} else if err != nil {
		return err
	}
	return err
}

func UpdateEventCategoryPriorities(m map[int]int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for key, val := range m {
		_, err := tx.Exec(`UPDATE event_categories SET priority=? WHERE id=?;`, val, key)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
