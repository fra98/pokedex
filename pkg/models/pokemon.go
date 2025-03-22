package models

// PokemonResponse represents a Pokemon response.
type PokemonResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"isLegendary"`
}
