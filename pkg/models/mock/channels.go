package mock

import (
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type ChannelModel struct{}

var TwitchChannel = &models.Channel{
	ID:              "42508152",
	PlatformID:      1,
	Login:           "starcraft",
	Title:           "StarCraft",
	ProfileImageURL: "https://static-cdn.jtvnw.net/jtv_user_pictures/0c9813cae3797d96-profile_image-300x300.png",
	Pattern:         "wcs",
	EventCategoryID: 1,
}

var YoutubeChannel = &models.Channel{
	ID:              "UCK5eBtuoj_HkdXKHNmBLAXg",
	PlatformID:      2,
	Login:           "",
	Title:           "AfreecaTV eSports",
	ProfileImageURL: "https://yt3.ggpht.com/a-/AAuE7mBZ1no98oeHv-OkWsyXSL7I9Fuj9LjPZ2JcHg=s88-mo-c-c0xffffffff-rj-k-no",
	Pattern:         "gsl",
	EventCategoryID: 6,
}

var Channels = []*models.Channel{
	{"42508152", 1, "starcraft", "StarCraft", "http://imageurl.com", "wcs", 1},
	{"UCK5eBtuoj_HkdXKHNmBLAXg", 2, "", "AfreecaTV eSports", "http://imageurl.com", "gsl", 2},
}

func (m *ChannelModel) SelectAllFromTwitch() ([]*models.Channel, error) {
	return []*models.Channel{TwitchChannel}, nil
}

func (m *ChannelModel) SelectFromAllCategories(platformID int) ([]*models.Channel, error) {
	switch platformID {
	case 1:
		return []*models.Channel{TwitchChannel}, nil
	default:
		return []*models.Channel{YoutubeChannel}, nil
	}
}

func (m *ChannelModel) SelectByCategory(categoryID int, platformID int) ([]*models.Channel, error) {
	fmt.Println(categoryID)
	fmt.Println(platformID)
	if categoryID == 1 && (platformID == 0 || platformID == 1) {
		return []*models.Channel{TwitchChannel}, nil
	} else if categoryID == 6 && (platformID == 0 || platformID == 2) {
		return []*models.Channel{YoutubeChannel}, nil
	} else {
		return []*models.Channel{}, nil
	}
}

func (m *ChannelModel) Insert(channel models.Channel, categoryID int) (*models.Channel, error) {
	channel.EventCategoryID = categoryID
	return &channel, nil
}

func (m *ChannelModel) DeleteFromCategory(channelID string, categoryID int) error {
	if channelID == "42508152" && categoryID == 1 {
		return nil
	} else if channelID == "UCK5eBtuoj_HkdXKHNmBLAXg" && categoryID == 2 {
		return nil
	} else {
		return models.ErrNotFound
	}
}
