package mock

import (
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

type ChannelModel struct{}

var Channels = []*models.Channel{
	{"42508152", 1, "starcraft", "StarCraft", "http://imageurl.com", "wcs", 1},
	{"UCK5eBtuoj_HkdXKHNmBLAXg", 2, "", "AfreecaTV eSports", "http://imageurl.com", "gsl", 2},
}

func GetPlatformChannels(platformID int) []*models.Channel {
	var channels []*models.Channel
	for _, c := range Channels {
		if c.PlatformID == platformID {
			channels = append(channels, c)
		}
	}
	return channels
}

func GetCategoryChannels(categoryID int, platformID int) []*models.Channel {
	channels := make([]*models.Channel, 0, 2)
	for _, c := range Channels {
		if c.EventCategoryID == categoryID && (platformID == 0 || platformID == c.PlatformID) {
			channels = append(channels, c)
		}
	}
	return channels
}

func (m *ChannelModel) SelectAllFromTwitch() ([]*models.Channel, error) {
	return GetPlatformChannels(1), nil
}

func (m *ChannelModel) SelectFromAllCategories(platformID int) ([]*models.Channel, error) {
	switch platformID {
	case 1:
		return GetPlatformChannels(1), nil
	case 2:
		return GetPlatformChannels(2), nil
	default:
		return Channels, nil
	}
}

func (m *ChannelModel) SelectByCategory(categoryID int, platformID int) ([]*models.Channel, error) {
	return GetCategoryChannels(categoryID, platformID), nil
}

func (m *ChannelModel) Insert(channel models.Channel, categoryID int) (*models.Channel, error) {
	if channel.ID == "UCK5eBtuoj_HkdXKHNmBLAXg" && categoryID == 2 {
		return nil, errors.New("error: Duplicate entry")
	}
	return &channel, nil
}

func (m *ChannelModel) DeleteFromCategory(channelID string, categoryID int) error {
	for _, c := range Channels {
		if c.ID == channelID && c.EventCategoryID == categoryID {
			return nil
		}
	}
	return models.ErrNotFound
}
