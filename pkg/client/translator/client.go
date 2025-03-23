package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/utils/ptr"

	"github.com/fra98/pokedex/pkg/consts"
	"github.com/fra98/pokedex/pkg/errors"
)

const defaultBaseURL = "https://api.funtranslations.com"

var _ Client = &FunTranslationClient{} // check if it implements the Client interface.

// FunTranslationClient represents a client that interacts with the FunTranslations API.
type FunTranslationClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewFunTranslationClient returns a new FunTranslations client.
func NewFunTranslationClient(baseURL *string) *FunTranslationClient {
	return &FunTranslationClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    ptr.Deref(baseURL, defaultBaseURL),
	}
}

type translationResponse struct {
	Contents struct {
		Translated string `json:"translated"`
	} `json:"contents"`
}

// Translate returns a translated text according to the translation type.
func (c *FunTranslationClient) Translate(ctx context.Context, text, translationType string) (string, error) {
	endpoint, err := c.getEndpoint(translationType)
	if err != nil {
		return "", err
	}

	requestBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limit exceeded error
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("failed to translate text: %w", errors.ErrRateLimitExceeded)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to translate text (code: %d): %w", resp.StatusCode, errors.ErrFailedRequest)
	}

	var translation translationResponse
	if err := json.NewDecoder(resp.Body).Decode(&translation); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return translation.Contents.Translated, nil
}

func (c *FunTranslationClient) getEndpoint(translationType string) (string, error) {
	switch translationType {
	case consts.YodaTranslationType:
		return c.baseURL + "/translate" + "/yoda.json", nil
	case consts.ShakespeareTranslationType:
		return c.baseURL + "/translate" + "/shakespeare.json", nil
	default:
		return "", fmt.Errorf("failed to retrieve endpoint: %w", errors.ErrUnsupportedTranslationType)
	}
}
