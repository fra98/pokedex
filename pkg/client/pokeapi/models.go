package pokeapi

// PokemonSpecies represents a Pokemon species.
type PokemonSpecies struct {
	Name              string            `json:"name"`
	IsLegendary       bool              `json:"is_legendary"`
	Habitat           Habitat           `json:"habitat"`
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
}

// Habitat represents a habitat where a Pokemon species can be found.
type Habitat struct {
	Name string `json:"name"`
}

// FlavorTextEntry represents a flavor text entry for a Pokemon species.
type FlavorTextEntry struct {
	FlavorText string   `json:"flavor_text"`
	Language   Language `json:"language"`
}

// Language represents a language.
type Language struct {
	Name string `json:"name"`
}
