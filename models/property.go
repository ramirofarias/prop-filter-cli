package models

type Property struct {
	SquareFootage float64         `json:"squareFootage"`
	Lighting      string          `json:"lighting"`
	Price         float64         `json:"price"`
	Rooms         float64         `json:"rooms"`
	Bathrooms     float64         `json:"bathrooms"`
	Location      [2]float64      `json:"location"`
	Description   string          `json:"description"`
	Ammenities    map[string]bool `json:"ammenities"`
}
