package translator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ Client = &CachedTranslationClient{} // check if it implements the Client interface.

// CachedTranslationClient represents a client that interacts with a translation API and caches the results.
type CachedTranslationClient struct {
	client Client
	cache  *cache.Cache
}

// NewCachedTranslationClient returns a new cached TranslationClient.
func NewCachedTranslationClient(client Client, timeoutExpiration, cleanupInterval time.Duration) *CachedTranslationClient {
	return &CachedTranslationClient{
		client: client,
		cache:  cache.New(timeoutExpiration, cleanupInterval),
	}
}

// Translate returns a translated text according to the translation type.
func (c *CachedTranslationClient) Translate(ctx context.Context, text, translationType string) (string, error) {
	cacheKey := "translation:" + translationType + ":" + text

	// Try to get from cache first
	if cachedData, found := c.cache.Get(cacheKey); found {
		cachedTranslation, ok := cachedData.(string)
		if ok {
			return cachedTranslation, nil
		}
		// Otherwise, remove the invalid cache entry and proceed
		log.Printf("Invalid cache entry for key %q", cacheKey)
		c.cache.Delete(cacheKey)
	}

	// Call the underlying client
	translation, err := c.client.Translate(ctx, text, translationType)
	if err != nil {
		return "", fmt.Errorf("failed to get translation: %w", err)
	}

	// Cache the result
	c.cache.Set(cacheKey, translation, 0)

	return translation, nil
}
