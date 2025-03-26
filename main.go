package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gitsang/configer"
	"github.com/spf13/cobra"
)

type Adapter struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

type Config struct {
	Adapter  string             `yaml:"adapter"`
	Adapters map[string]Adapter `yaml:"adapters"`
}

var rootCmd = &cobra.Command{
	Use: "goldendict-llm",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var rootFlags = struct {
	ConfigPaths []string
}{}

var cfger *configer.Configer

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&rootFlags.ConfigPaths, "config", "c", nil, "config file path")

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

func run() {
	var c Config
	err := cfger.Load(&c, rootFlags.ConfigPaths...)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	adapterConfig, ok := c.Adapters[c.Adapter]
	if !ok {
		panic(fmt.Errorf("adapter %s not found", c.Adapter))
	}

	reqBody := Request{
		Model: adapterConfig.Model,
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(fmt.Sprintf("Json marshal failed: %v", err))
	}

	req, err := http.NewRequest("POST", adapterConfig.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(fmt.Errorf("NewRequest failed: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adapterConfig.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("API request failed: %v", err))
	}
	defer resp.Body.Close()

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
		fmt.Println(apiResp.Choices[0].Message.Content)
	}
}

func main() {
	rootCmd.Execute()
}
