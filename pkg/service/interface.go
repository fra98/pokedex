package service

import (
	"context"

	"github.com/fra98/pokedex/pkg/models"
)

// Pokemon is an interface that defines the methods for retrieving Pokemon information.
type Pokemon interface {
	GetPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error)
	GetTranslatedPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error)
}
