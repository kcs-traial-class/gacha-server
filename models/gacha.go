package models

// GachaRequest represents the request body for performing gacha.
type GachaRequest struct {
	Times int `json:"times"`
}

// GachaResponse represents the response body after performing gacha.
type GachaResponse struct {
	Results []Item `json:"results"`
}
