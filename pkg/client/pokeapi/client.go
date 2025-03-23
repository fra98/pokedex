package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/utils/ptr"

	"github.com/fra98/pokedex/pkg/errors"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

var _ Client = &PokeAPIClient{} // check if it implements the Client interface.

// PokeAPIClient represents a client that interacts with the PokeAPI.
type PokeAPIClient struct { //nolint:revive // avoid conflict with the interface name
	httpClient *http.Client
	baseURL    string
}

// NewPokeAPIClient returns a new PokeAPIClient.
func NewPokeAPIClient(baseURL *string) *PokeAPIClient {
	return &PokeAPIClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    ptr.Deref(baseURL, defaultBaseURL),
	}
}

// GetPokemonSpecies returns a Pokemon species by name.
func (c *PokeAPIClient) GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error) {
	url := c.baseURL + "/pokemon-species/" + name

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get Pokemon species (code: %d): %w", resp.StatusCode, errors.ErrFailedRequest)
	}

	var species PokemonSpecies
	if err := json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &species, nil
}
