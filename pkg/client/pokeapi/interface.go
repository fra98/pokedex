package pokeapi

import "context"

// PokemonClient it's the interface that defines the methods to interact with the PokeAPI.
type PokemonClient interface {
	GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error)
}
