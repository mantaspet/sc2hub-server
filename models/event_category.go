package models

type EventCategory struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	InfoURL string `json:"info_url"`
	Order   int    `json:"order"`
}
