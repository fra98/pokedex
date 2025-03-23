package pokeapi

import "context"

// Client is an interface that defines the methods to retrieve Pokemon information from an API.
type Client interface {
	GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error)
}
