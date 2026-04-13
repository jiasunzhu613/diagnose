package api

import (
	"context"

	"google.golang.org/genai"
)

// Client encapsulates client state for interacting with the ollama
// service. Use [ClientFromEnvironment] to create new Clients.
type GoogleClient struct {
	sdk genai.Client
}

// Performs an HTTP request
func (c *GoogleClient) Do(ctx context.Context) {

}

func (c *GoogleClient) GenerateCompletion(ctx context.Context, req CompletionRequest) CompletionResponse {
	return CompletionResponse {}
}