package funtranslator

import "context"

// TranslatorClient it's the interface that defines the methods to interact with the funtranslations API.
type TranslatorClient interface {
	Translate(ctx context.Context, text, translationType string) (string, error)
}
