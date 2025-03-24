// Package integration_test provides integration tests for the application.
package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fra98/pokedex/pkg/api"
	"github.com/fra98/pokedex/pkg/client/pokeapi"
	"github.com/fra98/pokedex/pkg/client/translator"
	"github.com/fra98/pokedex/pkg/consts"
	"github.com/fra98/pokedex/pkg/models"
	"github.com/fra98/pokedex/pkg/service"
)

type testCase struct {
	name                    string
	pokemonName             string
	habitat                 string
	isLegendary             bool
	originalDescription     string
	yodaTranslation         string
	shakespeareTranslation  string
	expectedTranslationType string // "yoda", "shakespeare", or "original"
}

func TestPokemonAPIHandler(t *testing.T) { //nolint:funlen // skip long func length for test cases
	t.Parallel()

	// Test cases with different Pokemon types and expected translations
	testCases := []testCase{
		{
			name:                    "cave_pokemon",
			pokemonName:             "zubat",
			habitat:                 consts.HabitatCaveType,
			isLegendary:             false,
			originalDescription:     "Original description for cave Pokemon",
			yodaTranslation:         "Yoda translation for cave Pokemon, this is",
			shakespeareTranslation:  "Shakespeare translation for cave Pokemon",
			expectedTranslationType: consts.YodaTranslationType, // Should use Yoda for cave Pokemon
		},
		{
			name:                    "legendary_pokemon",
			pokemonName:             "mewtwo",
			habitat:                 "rare",
			isLegendary:             true,
			originalDescription:     "Original description for legendary Pokemon",
			yodaTranslation:         "Yoda translation for legendary Pokemon, this is",
			shakespeareTranslation:  "Shakespeare translation for legendary Pokemon",
			expectedTranslationType: consts.YodaTranslationType, // Should use Yoda for legendary Pokemon
		},
		{
			name:                    "cave_and_legendary_pokemon",
			pokemonName:             "registeel",
			habitat:                 consts.HabitatCaveType,
			isLegendary:             true,
			originalDescription:     "Original description for cave legendary Pokemon",
			yodaTranslation:         "Yoda translation for cave legendary Pokemon, this is",
			shakespeareTranslation:  "Shakespeare translation for cave legendary Pokemon",
			expectedTranslationType: consts.YodaTranslationType, // Should use Yoda for both cave and legendary
		},
		{
			name:                    "normal_pokemon",
			pokemonName:             "pikachu",
			habitat:                 "forest",
			isLegendary:             false,
			originalDescription:     "Original description for normal Pokemon",
			yodaTranslation:         "Yoda translation for normal Pokemon, this is",
			shakespeareTranslation:  "Shakespeare translation for normal Pokemon",
			expectedTranslationType: consts.ShakespeareTranslationType, // Should use Shakespeare for normal Pokemon
		},
		{
			name:                    "translation_failed",
			pokemonName:             "eevee",
			habitat:                 "urban",
			isLegendary:             false,
			originalDescription:     "Original description for failed translation",
			yodaTranslation:         "",         // Empty to simulate translation failure
			shakespeareTranslation:  "",         // Empty to simulate translation failure
			expectedTranslationType: "original", // Neither yoda or shakespeare, should fall back to original if translation fails
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Setup test server for PokeAPI that returns data based on test case
			pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Extract pokemon name from URL
				pokemonName := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

				// Verify it's the correct Pokemon for this test case
				assert.Equal(t, tc.pokemonName, pokemonName)

				// Return response based on test case
				response := fmt.Sprintf(`{
					"name": "%s",
					"is_legendary": %t,
					"habitat": {"name": "%s"},
					"flavor_text_entries": [
						{
							"flavor_text": "%s",
							"language": {"name": "en"}
						}
					]
				}`, tc.pokemonName, tc.isLegendary, tc.habitat, tc.originalDescription)

				_, err := w.Write([]byte(response))
				assert.NoError(t, err)
			}))
			defer pokeServer.Close()

			// Setup test server for Translator API
			translatorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				// Check which translation is being requested
				path := r.URL.Path
				var translation string

				switch {
				case strings.Contains(path, consts.YodaTranslationType):
					// Only return a valid translation if we have one for the test case
					if tc.yodaTranslation != "" {
						translation = tc.yodaTranslation
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusTooManyRequests) // Simulate rate limiting
						_, err := w.Write([]byte(`{"error":{"code":429,"message":"Too Many Requests"}}`))
						assert.NoError(t, err)
						return
					}
				case strings.Contains(path, consts.ShakespeareTranslationType):
					// Only return a valid translation if we have one for the test case
					if tc.shakespeareTranslation != "" {
						translation = tc.shakespeareTranslation
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusTooManyRequests) // Simulate rate limiting
						_, err := w.Write([]byte(`{"error":{"code":429,"message":"Too Many Requests"}}`))
						assert.NoError(t, err)
						return
					}
				default:
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				response := fmt.Sprintf(`{
					"contents": {
						"translated": "%s"
					}
				}`, translation)

				_, err := w.Write([]byte(response))
				assert.NoError(t, err)
			}))
			defer translatorServer.Close()

			// Create clients pointing to test servers
			pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)
			translatorClient := translator.NewFunTranslationClient(&translatorServer.URL)

			// Create service and handler
			pokemonService := service.NewPokemonService(pokeClient, translatorClient)
			pokemonHandler := api.NewPokemonHandler(pokemonService)

			// Set up router
			router := gin.New()
			router.GET("/pokemon/:name", pokemonHandler.GetPokemon)
			router.GET("/pokemon/translated/:name", pokemonHandler.GetTranslatedPokemon)

			// Test standard endpoint first
			w := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/pokemon/"+tc.pokemonName, http.NoBody)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var basicResponse models.PokemonResponse
			err = json.Unmarshal(w.Body.Bytes(), &basicResponse)
			require.NoError(t, err)

			assert.Equal(t, tc.pokemonName, basicResponse.Name)
			assert.Equal(t, tc.originalDescription, basicResponse.Description)
			assert.Equal(t, tc.habitat, basicResponse.Habitat)
			assert.Equal(t, tc.isLegendary, basicResponse.IsLegendary)

			// Test translated endpoint
			w = httptest.NewRecorder()
			req, err = http.NewRequestWithContext(t.Context(), http.MethodGet, "/pokemon/translated/"+tc.pokemonName, http.NoBody)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var translatedResponse models.PokemonResponse
			err = json.Unmarshal(w.Body.Bytes(), &translatedResponse)
			require.NoError(t, err)

			assert.Equal(t, tc.pokemonName, translatedResponse.Name)
			assert.Equal(t, tc.habitat, translatedResponse.Habitat)
			assert.Equal(t, tc.isLegendary, translatedResponse.IsLegendary)

			// Check that the correct translation was used
			switch tc.expectedTranslationType {
			case consts.YodaTranslationType:
				assert.Equal(t, tc.yodaTranslation, translatedResponse.Description)
			case consts.ShakespeareTranslationType:
				assert.Equal(t, tc.shakespeareTranslation, translatedResponse.Description)
			default:
				assert.Equal(t, tc.originalDescription, translatedResponse.Description)
			}
		})
	}
}
