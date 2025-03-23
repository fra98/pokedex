package pokeapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ Client = &CachedPokeAPIClient{} // check if it implements the Client interface.

// CachedPokeAPIClient represents a client that interacts with the PokeAPI and caches the results.
type CachedPokeAPIClient struct {
	client Client
	cache  *cache.Cache
}

// NewCachedPokeAPIClient returns a new cached PokeAPIClient.
func NewCachedPokeAPIClient(client Client, timeoutExpiration, cleanupInterval time.Duration) *CachedPokeAPIClient {
	return &CachedPokeAPIClient{
		client: client,
		cache:  cache.New(timeoutExpiration, cleanupInterval),
	}
}

// GetPokemonSpecies returns a Pokemon species by name.
func (c *CachedPokeAPIClient) GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error) {
	cacheKey := "pokeapi:species:" + name

	// Try to get from cache first
	if cachedData, found := c.cache.Get(cacheKey); found {
		cachedSpecies, ok := cachedData.(*PokemonSpecies)
		if ok {
			return cachedSpecies, nil
		}
		// Otherwise, remove the invalid cache entry and proceed
		log.Printf("Invalid cache entry for key %q", cacheKey)
		c.cache.Delete(cacheKey)
	}

	// Call the underlying client
	species, err := c.client.GetPokemonSpecies(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pokemon species: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, species, 0)

	return species, nil
}
