package models

var PlayerPageLength = 36

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

type PaginatedPlayers struct {
	Items  []*Player
	Cursor *int
}

type PlayerVideo struct {
	PlayerID int
	VideoID  string
}
