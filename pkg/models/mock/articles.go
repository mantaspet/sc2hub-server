package mock

import (
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"time"
)

func GetMockArticles() []*models.Article {
	time1 := time.Date(2019, 4, 17, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2019, 4, 9, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2019, 4, 5, 0, 0, 0, 0, time.UTC)
	articles := []*models.Article{
		{
			1,
			"Super(Hero) Tournament Preview, Part 1",
			"TeamLiquid.net",
			time1,
			"",
			"https://i.imgur.com/e2o9EmK.jpg",
			"https://tl.net/forum/starcraft-2/546128-superhero-tournament-preview-part-1",
		},
		{
			2,
			"WCS Spring: Tickets Now on Sale, Player Registrations Open!",
			"Blizzard",
			time2,
			"Information for players, community casters and spectators.",
			"https://bnetcmsus-a.akamaihd.net/cms/blog_thumbnail/4i/4IUKYJR5JREV1554503304373.jpg",
			"https://news.blizzard.com/en-us/blizzard/22945355/wcs-spring-tickets-now-on-sale-player-registrations-open",
		},
		{
			3,
			"WCS Winter Americas: Playoffs Preview",
			"TeamLiquid.net",
			time3,
			"",
			"https://i.imgur.com/e2o9EmK.jpg",
			"https://tl.net/forum/starcraft-2/545441-wcs-winter-americas-playoffs-preview",
		},
	}
	return articles
}
