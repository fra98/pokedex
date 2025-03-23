package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/fra98/pokedex/pkg/client/pokeapi"
	"github.com/fra98/pokedex/pkg/client/translator"
	"github.com/fra98/pokedex/pkg/consts"
	"github.com/fra98/pokedex/pkg/errors"
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
func (s *PokemonService) GetPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error) {
	pokemonSpecies, err := s.pokeClient.GetPokemonSpecies(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve pokemon species: %w", err)
	}

	description, err := extractEnglishDescription(pokemonSpecies)
	if err != nil {
		return nil, fmt.Errorf("unable to extract English description: %w", err)
	}

	return &models.PokemonResponse{
		Name:        pokemonSpecies.Name,
		Description: description,
		Habitat:     pokemonSpecies.Habitat.Name,
		IsLegendary: pokemonSpecies.IsLegendary,
	}, nil
}

// GetTranslatedPokemonInfo retrieves the information of a Pokemon given its name with a translated description.
func (s *PokemonService) GetTranslatedPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error) {
	// Get basic info first
	pokemon, err := s.GetPokemonInfo(ctx, name)
	if err != nil {
		return nil, err
	}

	// Determine translation type
	var translationType string
	if pokemon.IsLegendary || pokemon.Habitat == consts.HabitatCaveType {
		translationType = consts.YodaTranslationType
	} else {
		translationType = consts.ShakespeareTranslationType
	}

	// Get translation
	translatedDesc, err := s.translatorClient.Translate(ctx, pokemon.Description, translationType)
	if err != nil {
		// if translation fails, fallback to original description
		translatedDesc = pokemon.Description
	}

	// Update description
	pokemon.Description = translatedDesc
	return pokemon, nil
}

// Helper function to extract English description.
func extractEnglishDescription(species *pokeapi.PokemonSpecies) (string, error) {
	for i := range species.FlavorTextEntries {
		entryFlavor := &species.FlavorTextEntries[i]
		if entryFlavor.Language.Name == "en" {
			return sanitizeDescription(entryFlavor.FlavorText), nil
		}
	}
	return "", errors.ErrResourceNotFound
}

// Helper function to sanitize description.
func sanitizeDescription(description string) string {
	sanitizedDesc := strings.ReplaceAll(description, "\n", " ")  // replace newlines with spaces
	sanitizedDesc = strings.ReplaceAll(sanitizedDesc, "\f", " ") // replace form feeds with spaces
	return sanitizedDesc
}
