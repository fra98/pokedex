package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fra98/pokedex/pkg/client/pokeapi"
	"github.com/fra98/pokedex/pkg/client/translator"
	"github.com/fra98/pokedex/pkg/service"
)

const (
	testPokemonLegendary = `{
		"name": "mewtwo",
		"is_legendary": true,
		"habitat": {"name": "rare"},
		"flavor_text_entries": [
			{
				"flavor_text": "original description",
				"language": {"name": "en"}
			}
		]
	}`

	testPokemonNotRareNotCave = `{
		"name": "pikachu",
		"is_legendary": false,
		"habitat": {"name": "forest"},
		"flavor_text_entries": [
			{
				"flavor_text": "original description",
				"language": {"name": "en"}
			}
		]
	}`

	testPokemonNoText = `{
		"name": "mewtwo",
		"is_legendary": true,
		"habitat": {"name": "rare"},
		"flavor_text_entries": [
			{
				"flavor_text": "description in other language",
				"language": {"name": "fr"}
			}
		]
	}`

	testContentTranslated = `{
		"contents": {
			"translated": "translated description"
		}
	}`
)

func TestGetPokemonInfo_Success(t *testing.T) {
	t.Parallel()

	// Create a new test server with the handler
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "/pokemon-species/mewtwo", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	}))
	defer pokeServer.Close()

	// Create a real client but point it to our test server
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)

	// We don't need the translator for this test
	translatorClient := translator.NewFunTranslationClient(nil) // Empty URL since we're not using it

	// Create the service with our stubbed clients
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "original description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetPokemonInfo_Failure(t *testing.T) {
	t.Parallel()

	// Create a new test server with the handler
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return an error
		w.WriteHeader(http.StatusNotFound)
	}))
	defer pokeServer.Close()

	// Create a real client but point it to our test server
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)

	// We don't need the translator for this test
	translatorClient := translator.NewFunTranslationClient(nil) // Empty URL since we're not using it

	// Create the service with our stubbed clients
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTranslatedPokemonInfo_Success(t *testing.T) {
	t.Parallel()

	// Setup test server for PokeAPI
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response for pokemon-species
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	}))
	defer pokeServer.Close()

	// Setup test server for Translator API
	translatorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/translate/yoda.json", r.URL.Path) // Should use Yoda for legendary

		// Parse request body
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)
		assert.Equal(t, "original description", requestBody["text"])

		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(testContentTranslated))
		assert.NoError(t, err)
	}))
	defer translatorServer.Close()

	// Create clients with test servers
	// Create a real client but point it to our test server
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)

	// We don't need the translator for this test
	translatorClient := translator.NewFunTranslationClient(&translatorServer.URL)

	// Create the service with our stubbed clients
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "translated description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_NoText(t *testing.T) {
	t.Parallel()

	// Setup test server for PokeAPI
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response for pokemon-species
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonNoText))
		assert.NoError(t, err)
	}))
	defer pokeServer.Close()

	// Create clients with test servers
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)
	translatorClient := translator.NewFunTranslationClient(nil) // Empty URL since we're not using it

	// Create the service
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get an error since there's no English text
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTranslatedPokemonInfo_RateLimiting(t *testing.T) {
	t.Parallel()

	// Setup PokeAPI test server with normal response
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	}))
	defer pokeServer.Close()

	// Setup Translator API with rate limit response
	translatorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return rate limiting error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, err := w.Write([]byte(`{"error": {"code": 429, "message": "Too Many Requests: Rate limit exceeded"}}`))
		assert.NoError(t, err)
	}))
	defer translatorServer.Close()

	// Create clients with test servers
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)
	translatorClient := translator.NewFunTranslationClient(&translatorServer.URL)

	// Create the service
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get the basic description back when translation fails
	require.NoError(t, err) // This should not return an error
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "original description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_Shakespeare(t *testing.T) {
	t.Parallel()

	// Setup test server for PokeAPI
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response for pokemon-species
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonNotRareNotCave))
		assert.NoError(t, err)
	}))
	defer pokeServer.Close()

	// Setup test server for Translator API
	translatorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/translate/shakespeare.json", r.URL.Path) // Should use Shakespeare for normal Pokemon

		// Parse request body
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)
		assert.Equal(t, "original description", requestBody["text"])

		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(testContentTranslated))
		assert.NoError(t, err)
	}))
	defer translatorServer.Close()

	// Create clients with test servers
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)
	translatorClient := translator.NewFunTranslationClient(&translatorServer.URL)

	// Create the service
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "pikachu")

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, "pikachu", result.Name)
	assert.Equal(t, "translated description", result.Description)
	assert.Equal(t, "forest", result.Habitat)
	assert.False(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_FailurePoke(t *testing.T) {
	t.Parallel()

	// Setup PokeAPI test server with error response
	pokeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer pokeServer.Close()

	// Create clients with test servers
	pokeClient := pokeapi.NewPokeAPIClient(&pokeServer.URL)
	translatorClient := translator.NewFunTranslationClient(nil) // Empty URL since we're not using it

	// Create the service
	pokemonService := service.NewPokemonService(pokeClient, translatorClient)

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get an error
	require.Error(t, err)
	assert.Nil(t, result)
}
