package service

import "github.com/fra98/pokedex/pkg/models"

// Pokemon is an interface that defines the methods for retrieving Pokemon information.
type Pokemon interface {
	GetPokemonInfo(name string) (*models.PokemonResponse, error)
	GetTranslatedPokemonInfo(name string) (*models.PokemonResponse, error)
}
