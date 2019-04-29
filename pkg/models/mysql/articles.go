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

func (m *ArticleModel) SelectPage(fromDate string, query string) ([]*models.Article, error) {
	words := strings.Fields(query)
	valueArgs := make([]interface{}, 0, len(words)*2+2)
	stmt := `SELECT
			id, title, source, published_at, COALESCE(excerpt, ''), thumbnail_url, url
	  	FROM articles`

	if fromDate == "" {
		stmt += " WHERE 1<>?"
	} else {
		stmt += " WHERE published_at<=?"
	}

	valueArgs = append(valueArgs, fromDate)
	for _, w := range words {
		stmt += ` AND (title LIKE ? OR source LIKE ?)`
		valueArgs = append(valueArgs, "%"+w+"%")
		valueArgs = append(valueArgs, "%"+w+"%")
	}
	valueArgs = append(valueArgs, models.ArticlePageLength+1)
	stmt += " ORDER BY published_at DESC LIMIT ?"

	rows, err := m.DB.Query(stmt, valueArgs...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*models.Article{}
	for rows.Next() {
		a := &models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Source, &a.PublishedAt, &a.Excerpt, &a.ThumbnailURL, &a.URL)
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

func (m *ArticleModel) SelectRecent() ([]*models.Article, error) {
	stmt := `SELECT id, title, source, published_at, COALESCE(excerpt, ''), thumbnail_url, url
	  	FROM articles
	  	ORDER BY published_at DESC
	  	LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*models.Article{}
	for rows.Next() {
		a := &models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Source, &a.PublishedAt, &a.Excerpt, &a.ThumbnailURL, &a.URL)
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

func (m *ArticleModel) SelectByCategory(categoryID int) ([]models.Article, error) {
	var articles []models.Article
	//articles := []models.Article{
	//	models.Article{
	//		1, 1, categoryID, 1, "Maru wins yet another GSL", "John Doe", "2019-04-01", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
	//	},
	//	models.Article{
	//		2, 1, categoryID, 1, "Koreans continue to dominate", "A sad foreigner", "2019-03-25", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
	//	},
	//	models.Article{
	//		3, 1, categoryID, 1, "What does it take to get to silver league in Starcraft 2", "A platinum veteran", "2019-04-01", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur posuere elit et ligula cursus, sed fermentum erat faucibus. Morbi sollicitudin et quam sagittis mattis. Donec pellentesque, lorem id congue interdum, lorem tellus gravida diam, in mattis lorem quam ac erat. Sed enim est, condimentum eu pretium eget, pretium eu massa. Cras vestibulum porttitor tellus, vel ullamcorper dui scelerisque ultrices. Aliquam ultrices justo ligula, quis elementum dui interdum non. Nullam sed viverra neque. In hac habitasse platea dictumst. Vestibulum vel justo non erat condimentum ornare ut ac ante. Donec condimentum lobortis risus in facilisis. Aliquam imperdiet tellus ut lectus varius, non tincidunt ipsum vehicula. Vestibulum nec arcu non leo dapibus semper id vel risus.\n\nPhasellus elementum aliquet nisi, eu fermentum diam commodo sit amet. Etiam eu tristique risus. Suspendisse condimentum finibus congue. Etiam porttitor lorem id massa auctor aliquet. Pellentesque congue porta purus at rhoncus. Suspendisse nec ipsum et sem dignissim finibus ac at est. Donec sed ex odio. Duis quis interdum metus. Morbi ac ultricies neque. Mauris aliquam velit nec ligula mollis, et condimentum est dapibus. Mauris tincidunt odio at malesuada mattis. Praesent porta iaculis lectus, non ultrices risus vehicula sed. Donec quis dolor felis. Ut lacus tellus, suscipit nec mauris vel, laoreet iaculis libero. Pellentesque a nunc quis ligula tristique dapibus. Nunc velit diam, congue vitae feugiat vitae, congue ac elit.", "https://deepmind.com/blog/deepmind-and-blizzard-open-starcraft-ii-ai-research-environment/",
	//	},
	//}

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
