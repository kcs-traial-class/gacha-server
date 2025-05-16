package models

type Item struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rarity          string  `json:"rarity"`
	Details         string  `json:"details"`
	Percentage      float64 `json:"percentage"`
	ImageIdentifier string  `json:"image_identifier"`
}
