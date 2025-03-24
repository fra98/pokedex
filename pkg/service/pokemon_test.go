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

// setupService is a helper function to setup the test servers for the poke and translator API.
// It accepts:
// - the test server handler function for the poke client
// - the test server handler function for the translator client
// If the handler is nil, the server is not created.
// It returns the test servers and the PokemonService with the clients configured to use the test servers.
func setupService(pokeHandler, translatorHandler *http.HandlerFunc) (pokeServer, translServer *httptest.Server, pokeService *service.PokemonService) {
	// Setup test server for PokeAPI
	if pokeHandler != nil {
		pokeServer = httptest.NewServer(*pokeHandler)
	}

	if translatorHandler != nil {
		translServer = httptest.NewServer(*translatorHandler)
	}

	// Create clients with test servers
	var pokeClient pokeapi.Client
	if pokeServer != nil {
		pokeClient = pokeapi.NewPokeAPIClient(&pokeServer.URL)
	} else {
		pokeClient = pokeapi.NewPokeAPIClient(nil)
	}

	var translatorClient translator.Client
	if translServer != nil {
		translatorClient = translator.NewFunTranslationClient(&translServer.URL)
	} else {
		translatorClient = translator.NewFunTranslationClient(nil)
	}

	// Create the service
	pokeService = service.NewPokemonService(pokeClient, translatorClient)

	return pokeServer, translServer, pokeService
}

func TestGetPokemonInfo_Success(t *testing.T) {
	t.Parallel()

	// Poke handler: return a successful response
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "/pokemon-species/mewtwo", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	})

	// Setup test servers and service
	pokeServer, _, pokemonService := setupService(&pokeHandler, nil) // no translator needed
	defer pokeServer.Close()

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get the original description
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "original description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetPokemonInfo_Failure(t *testing.T) {
	t.Parallel()

	// Poke handler: return an error
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return an error
		w.WriteHeader(http.StatusNotFound)
	})

	// Setup test servers and service
	pokeServer, _, pokemonService := setupService(&pokeHandler, nil) // no translator needed
	defer pokeServer.Close()

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTranslatedPokemonInfo_Success(t *testing.T) {
	t.Parallel()

	// Pokehandler: test legendary pokemon with Yoda translation
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	})

	// Translator handler: translate to Yoda
	translHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Setup test servers and service
	pokeServer, translServer, pokemonService := setupService(&pokeHandler, &translHandler)
	defer pokeServer.Close()
	defer translServer.Close()

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get a translated translation
	require.NoError(t, err)
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "translated description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_SuccessShakespeare(t *testing.T) {
	t.Parallel()

	// Poke handler: test pokemon with shakespeare translation (not legendary, not cave)
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonNotRareNotCave))
		assert.NoError(t, err)
	})

	// Translator handler: translate to Shakespeare
	translHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Setup test servers and service
	pokeServer, translServer, pokemonService := setupService(&pokeHandler, &translHandler)
	defer pokeServer.Close()
	defer translServer.Close()

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "pikachu")

	// Assertions - should get a translated translation
	require.NoError(t, err)
	assert.Equal(t, "pikachu", result.Name)
	assert.Equal(t, "translated description", result.Description)
	assert.Equal(t, "forest", result.Habitat)
	assert.False(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_NoText(t *testing.T) {
	t.Parallel()

	// Poke handler: test pokemon with no English text available
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return a fixed response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonNoText))
		assert.NoError(t, err)
	})

	// Setup test servers and service
	pokeServer, _, pokemonService := setupService(&pokeHandler, nil) // no translator needed
	defer pokeServer.Close()

	// Call the service method
	result, err := pokemonService.GetPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get an error since there's no English text
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTranslatedPokemonInfo_RateLimiting(t *testing.T) {
	t.Parallel()

	// Poke handler: successful response
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testPokemonLegendary))
		assert.NoError(t, err)
	})

	// Translator handler: rate limit response
	translHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Return rate limiting error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, err := w.Write([]byte(`{"error": {"code": 429, "message": "Too Many Requests: Rate limit exceeded"}}`))
		assert.NoError(t, err)
	})

	// Setup test servers and service
	pokeServer, translServer, pokemonService := setupService(&pokeHandler, &translHandler)
	defer pokeServer.Close()
	defer translServer.Close()

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get the original description back when translation fails
	require.NoError(t, err) // This should not return an error
	assert.Equal(t, "mewtwo", result.Name)
	assert.Equal(t, "original description", result.Description)
	assert.Equal(t, "rare", result.Habitat)
	assert.True(t, result.IsLegendary)
}

func TestGetTranslatedPokemonInfo_FailurePoke(t *testing.T) {
	t.Parallel()

	// Poke handler: returns an error
	pokeHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Setup test servers and service
	pokeServer, _, pokemonService := setupService(&pokeHandler, nil)
	defer pokeServer.Close()

	// Call the service method
	result, err := pokemonService.GetTranslatedPokemonInfo(t.Context(), "mewtwo")

	// Assertions - should get an error
	require.Error(t, err)
	assert.Nil(t, result)
}
