package service

import (
	"github.com/fra98/pokedex/pkg/models"
)

// PokemonService is an interface that defines the methods for retrieving Pokemon information.
type PokemonService interface {
	GetPokemonInfo(name string) (*models.PokemonResponse, error)
	GetTranslatedPokemonInfo(name string) (*models.PokemonResponse, error)
}
