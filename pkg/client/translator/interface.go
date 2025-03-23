package translator

import "context"

// Client is an interface that defines the methods to translate text from an API.
type Client interface {
	Translate(ctx context.Context, text, translationType string) (string, error)
}
