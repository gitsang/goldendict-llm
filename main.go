package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gitsang/configer"
	"github.com/spf13/cobra"
)

type Adapter struct {
	URL   string
	Token string
	Model string
}

type Config struct {
	Adapter  string
	Adapters map[string]Adapter
	Timeout  string `default:"30s"`
}

var rootCmd = &cobra.Command{
	Use: "goldendict-llm",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var rootFlags = struct {
	ConfigPaths []string
	UserInput   string
}{}

var cfger *configer.Configer

//go:embed static/word-prompt.md
var WordPrompt string

//go:embed static/sentence-prompt.md
var SentencePrompt string

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&rootFlags.ConfigPaths, "config", "c", nil, "config file path")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.UserInput, "content", "m", "", "user input content")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			rootFlags.UserInput = joinArgs(args)
		}
		run()
	}

	cfger = configer.New(
		configer.WithTemplate(new(Config)),
		configer.WithEnvBind(
			configer.WithEnvPrefix("GOLDEN_DICT_LLM"),
			configer.WithEnvDelim("_"),
		),
		configer.WithFlagBind(
			configer.WithCommand(rootCmd),
			configer.WithFlagPrefix(""),
			configer.WithFlagDelim("."),
		),
	)
}

func isWord(content string) bool {
	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 {
		return false
	}
	if strings.Contains(trimmed, " ") {
		return false
	}
	for _, sep := range []string{",", ".", ";", ":", "，", "。", "；", "：", "、"} {
		if strings.Contains(trimmed, sep) {
			return false
		}
	}

	chineseCharCount := 0
	englishCharCount := 0
	for _, r := range trimmed {
		if unicode.Is(unicode.Han, r) {
			chineseCharCount++
		} else if unicode.IsLetter(r) {
			englishCharCount++
		}

		if chineseCharCount > 4 {
			return false
		}
		if chineseCharCount == 0 && englishCharCount > 32 {
			return false
		}

		if chineseCharCount > 0 && (chineseCharCount+englishCharCount) > 16 {
			return false
		}
	}

	return true
}

func run() {
	var c Config
	err := cfger.Load(&c, rootFlags.ConfigPaths...)
	if err != nil {
		panic(err)
	}

	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid timeout: %v", err))
	}

	client := &http.Client{
		Timeout: timeout,
	}

	adapterConfig, ok := c.Adapters[c.Adapter]
	if !ok {
		panic(fmt.Errorf("adapter %s not found", c.Adapter))
	}

	promptTemplate := WordPrompt
	if !isWord(rootFlags.UserInput) {
		promptTemplate = SentencePrompt
	}

	renderedUserInput, err := RenderUserInputTemplateToString(rootFlags.UserInput)
	if err != nil {
		panic(fmt.Errorf("RenderUserInputTemplateToString failed: %v", err))
	}

	reqBody := Request{
		Model: adapterConfig.Model,
		Messages: []Message{
			{Role: "system", Content: promptTemplate},
			{Role: "user", Content: renderedUserInput},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(fmt.Sprintf("Json marshal failed: %v", err))
	}

	req, err := http.NewRequest("POST", adapterConfig.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(fmt.Errorf("NewRequest failed: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adapterConfig.Token)

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("API request failed: %v", err))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("ReadAll failed: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("API request failed with status: %d", resp.StatusCode))
	}

	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		panic(fmt.Sprintf("Json unmarshal failed: %v", err))
	}

	if len(apiResp.Choices) > 0 {
		content := apiResp.Choices[0].Message.Content
		if isWord(rootFlags.UserInput) {
			renderedContent, err := ProcessWordResponse(content)
			if err != nil {
				panic(fmt.Sprintf("Template rendering failed: %v", err))
			}
			fmt.Println(renderedContent)
		} else {
			renderedContent, err := RenderSentenceTemplateToString(rootFlags.UserInput, content)
			if err != nil {
				panic(fmt.Sprintf("Template rendering failed: %v", err))
			}
			fmt.Println(renderedContent)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
