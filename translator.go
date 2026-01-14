package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//go:embed static/word-prompt.md
var WordPrompt string

//go:embed static/sentence-prompt.md
var SentencePrompt string

type AdapterConfig struct {
	Name  string
	URL   string
	Token string
	Model string
}

// Translator handles translation requests with caching
type Translator struct {
	adapter AdapterConfig
	client  *http.Client
}

// TranslatorOption defines a function that configures a Translator
type TranslatorOption func(*Translator)

// WithHTTPClient sets a custom HTTP client for the translator
func WithHTTPClient(client *http.Client) TranslatorOption {
	return func(t *Translator) {
		t.client = client
	}
}

// NewTranslator creates a new Translator with the given options
func NewTranslator(adapter AdapterConfig, opts ...TranslatorOption) *Translator {
	translator := &Translator{
		adapter: adapter,
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt(translator)
	}

	return translator
}

// TranslateWord translates a word and returns the formatted HTML
func (t *Translator) TranslateWord(word string) (string, error) {
	var (
		startTime = time.Now()
		err       error
		result    string
		duration  string
	)

	renderedUserInput, err := RenderUserInputTemplateToString(word)
	if err != nil {
		return "", fmt.Errorf("RenderUserInputTemplateToString failed: %v", err)
	}

	reqBody := Request{
		Model: t.adapter.Model,
		Messages: []Message{
			{Role: "system", Content: WordPrompt},
			{Role: "user", Content: renderedUserInput},
		},
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %v", err)
	}

	req, err := http.NewRequest("POST", t.adapter.URL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", fmt.Errorf("NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.adapter.Token)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("readAll failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var responseBody Response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return "", fmt.Errorf("json unmarshal failed: %v", err)
	}

	if len(responseBody.Choices) > 0 {
		result = responseBody.Choices[0].Message.Content
	}
	duration = fmt.Sprintf("%.2fs", time.Since(startTime).Seconds())

	renderedContent, err := ProcessWordResponseWithAdapterInfo(result, t.adapter.Name, t.adapter.Model, duration)
	if err != nil {
		return "", err
	}
	return renderedContent, nil
}
