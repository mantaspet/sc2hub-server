package models

type Player struct {
	ID            int
	PlayerID      string
	Name          string
	Race          string
	Team          string
	Country       string
	TotalEarnings float32
	DateOfBirth   string
	LiquipediaURL string
	ImageURL      string
	StreamURL     string
	IsRetired     bool
}

type PlayerVideo struct {
	PlayerID int
	VideoID  string
}
