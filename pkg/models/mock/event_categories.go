package mock

import "github.com/mantaspet/sc2hub-server/pkg/models"

var EventCategories = []*models.EventCategory{
	{1, "World Championship Series", "wcs", "", "https://infourl.com", "http://imageurl.com", "", 1},
	{2, "Global Starcraft League Code S", "gsl", "", "https://infourl.com", "http://imageurl.com", "description", 2},
	{3, "Intel Extreme Masters", "iem", "", "https://infourl.com", "http://imageurl.com", "", 3},
}

var EventCategoryPatterns = []*models.EventCategory{
	{ID: 1, IncludePatterns: "wcs"},
	{ID: 2, IncludePatterns: "gsl"},
	{ID: 3, IncludePatterns: "iem"},
}

var EventCategoryArticles = []models.EventCategoryArticle{
	{1, 2},
	{1, 3},
}

var InvalidEventCategoryArticles = []models.EventCategoryArticle{
	{4, 2},
	{1, 5},
}
