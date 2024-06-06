package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	client        *http.Client
	systemMessage string
	model         string
	apiKey        string
}

func NewAPI(client *http.Client, systemMessage, model, apiKey string) *API {
	return &API{client: client, systemMessage: systemMessage, model: model, apiKey: apiKey}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completion struct {
	Model    string     `json:"model"`
	Messages []*message `json:"messages"`
}

type choice struct {
	Index   int      `json:"index"`
	Message *message `json:"message"`
}

type completionObject struct {
	Choices []*choice `json:"choices"`
}

func (a *API) Completions(ctx context.Context, q string) (answer string, err error) {
	messages := []*message{
		{Role: "system", Content: a.systemMessage},
		{Role: "user", Content: q},
	}
	completion := &completion{Model: a.model, Messages: messages}
	bodyBytes, err := json.Marshal(completion)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.apiKey)
	response, err := a.client.Do(request)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_err := Body.Close()
		if _err != nil {
			err = _err
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(`openai completions api non ok response: %s`, response.Status)
	}
	responseBytes, err := io.ReadAll(response.Body)
	completionObj := &completionObject{}
	err = json.Unmarshal(responseBytes, completionObj)
	if err != nil {
		return "", err
	}
	if len(completionObj.Choices) == 0 {
		return "", fmt.Errorf(`0 choices returned from openai completions endpoint`)
	}
	message := completionObj.Choices[0].Message
	return message.Content, nil
}
