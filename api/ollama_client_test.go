package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/jiasunzhu613/diagnose/envconfig"
)

func TestOllamaClient(t *testing.T) {
	var err error

	client, err := NewOllamaClient(envconfig.BASE_OLLAMA_LOCAL)
	if err != nil {
		t.Error("Failed to instantiate OllamaClient with BASE_OLLAMA_LOCAL path")
	}

	ctx := context.Background()
	request := &GenerateRequest{
		Model: envconfig.OLLAMA_QWEN,
		Prompt: "Why is the sky blue?",
	}

	handler := func(resp GenerateResponse) error {
		fmt.Print(resp.Response)

		return nil
	}

	err = client.GenerateCompletion(ctx, request, handler)
	if err != nil {
		t.Errorf("GenerateComplete failed, %v", err)
	}
}