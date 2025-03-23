package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/utils/ptr"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

// Client represents a PokeAPI client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient returns a new PokeAPI client.
func NewClient(baseURL *string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    ptr.Deref(baseURL, defaultBaseURL),
	}
}

// GetPokemonSpecies returns a Pokemon species by name.
func (c *Client) GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error) {
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
		return nil, fmt.Errorf("failed to get Pokemon species (code: %d): %w", resp.StatusCode, ErrFailedRequest)
	}

	var species PokemonSpecies
	if err := json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &species, nil
}
