package service

import (
	"github.com/fra98/pokedex/pkg/client/pokeapi"
	"github.com/fra98/pokedex/pkg/client/translator"
	"github.com/fra98/pokedex/pkg/models"
)

var _ Pokemon = &PokemonService{}

// PokemonService implements the Pokemon interface according to the API requirements.
type PokemonService struct {
	pokeClient       pokeapi.Client
	translatorClient translator.Client
}

// NewPokemonService creates a new PokemonService with the given clients.
func NewPokemonService(pokeClient pokeapi.Client, translatorClient translator.Client) *PokemonService {
	return &PokemonService{
		pokeClient:       pokeClient,
		translatorClient: translatorClient,
	}
}

// GetPokemonInfo retrieves the information of a Pokemon given its name.
func (s *PokemonService) GetPokemonInfo(name string) (*models.PokemonResponse, error) {
	_ = name
	panic("not implemented") // TODO: Implement
}

// GetTranslatedPokemonInfo retrieves the information of a Pokemon given its name with a translated description.
func (s *PokemonService) GetTranslatedPokemonInfo(name string) (*models.PokemonResponse, error) {
	_ = name
	panic("not implemented") // TODO: Implement
}
