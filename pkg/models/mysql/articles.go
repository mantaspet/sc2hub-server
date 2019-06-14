package mysql

import (
	"database/sql"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"strings"
)

type ArticleModel struct {
	DB *sql.DB
}

func parseArticleRows(rows *sql.Rows) ([]*models.Article, error) {
	articles := []*models.Article{}
	for rows.Next() {
		a := &models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Source, &a.PublishedAt, &a.Excerpt, &a.ThumbnailURL, &a.URL)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func queryArticlesPage(db *sql.DB, stmt string, pivotID int, pageSize int, from int, query string) ([]*models.Article, error) {
	words := strings.Fields(query)
	valueArgs := make([]interface{}, 0, len(words)*2+3)
	if pivotID > 0 {
		valueArgs = append(valueArgs, pivotID)
	}
	for _, w := range words {
		stmt += ` AND (title LIKE ? OR source LIKE ?)`
		valueArgs = append(valueArgs, "%"+w+"%")
		valueArgs = append(valueArgs, "%"+w+"%")
	}
	stmt += " ORDER BY articles.published_at DESC LIMIT ?,?"
	valueArgs = append(valueArgs, from, pageSize+1)

	rows, err := db.Query(stmt, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseArticleRows(rows)
}

func (m *ArticleModel) SelectPage(pageSize int, from int, query string) ([]*models.Article, error) {
	stmt := `SELECT id, title, source, published_at, COALESCE(excerpt, ''), thumbnail_url, url
	  	FROM articles
		WHERE 1=1`

	return queryArticlesPage(m.DB, stmt, 0, pageSize, from, query)
}

func (m *ArticleModel) SelectByPlayer(pageSize int, from int, query string, playerID int) ([]*models.Article, error) {
	stmt := `SELECT articles.id, articles.title, articles.source, articles.published_at, COALESCE(articles.excerpt, ''),
       	articles.thumbnail_url, articles.url
	  	FROM articles
	  	INNER JOIN player_articles pa
	  	ON articles.id = pa.article_id
	  	WHERE pa.player_id=?`

	return queryArticlesPage(m.DB, stmt, playerID, pageSize, from, query)
}

func (m *ArticleModel) SelectByCategory(pageSize int, from int, query string, categoryID int) ([]*models.Article, error) {
	stmt := `SELECT articles.id, articles.title, articles.source, articles.published_at, COALESCE(articles.excerpt, ''),
       	articles.thumbnail_url, articles.url
	  	FROM articles
	  	INNER JOIN event_category_articles eca
	  	ON articles.id = eca.article_id
	  	WHERE eca.event_category_id=?`

	return queryArticlesPage(m.DB, stmt, categoryID, pageSize, from, query)
}

func (m *ArticleModel) SelectLastInserted(amount int64) ([]*models.Article, error) {
	stmt := `
		SELECT id, title, COALESCE(excerpt, '')
		FROM articles
		ORDER BY id DESC
		LIMIT ?`

	rows, err := m.DB.Query(stmt, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*models.Article{}
	for rows.Next() {
		a := &models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Excerpt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (m *ArticleModel) InsertMany(articles []models.Article) (int64, error) {
	valueStrings := make([]string, 0, len(articles))
	valueArgs := make([]interface{}, 0, len(articles)*6)
	for _, e := range articles {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, e.Title)
		valueArgs = append(valueArgs, e.Source)
		valueArgs = append(valueArgs, e.PublishedAt)
		valueArgs = append(valueArgs, e.Excerpt)
		valueArgs = append(valueArgs, e.ThumbnailURL)
		valueArgs = append(valueArgs, e.URL)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO articles(title, source, published_at, excerpt, thumbnail_url, url)
		VALUES %s 
		ON DUPLICATE KEY UPDATE
			thumbnail_url=VALUES(thumbnail_url),
			url=VALUES(url),
			excerpt=VALUES(excerpt);`, strings.Join(valueStrings, ","))

	res, err := m.DB.Exec(stmt, valueArgs...)
	_, _ = m.DB.Exec(`ALTER TABLE articles AUTO_INCREMENT=1`) // to prevent ON DUPLICATE KEY triggers from inflating next ID
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowCnt, nil
}
