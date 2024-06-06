package chatgpt

import (
	"context"
	"github.com/zain-saqer/twitch-chatgpt/internal/env"
	"net/http"
	"testing"
)

func TestApiChatCompletion(t *testing.T) {
	apiKey, err := env.GetEnv(`OPENAI_API_KEY`)
	if err != nil {
		t.Fatal(err)
	}
	api := NewAPI(&http.Client{}, `repeat exactly what the user say`, `gpt-3.5-turbo`, apiKey)
	q := `ping 123`
	answer, err := api.Completions(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	if answer != q {
		t.Errorf("got %q, want %q", answer, q)
	}
}
