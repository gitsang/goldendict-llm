package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gitsang/configer"
	"github.com/spf13/cobra"
)

type Config struct {
	Adapter  string
	Adapters map[string]AdapterConfig
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
			configer.WithEnvPrefix("GOLDENDICT_LLM"),
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

	// adapter
	adapterConfig, ok := c.Adapters[c.Adapter]
	if !ok {
		panic(fmt.Errorf("adapter %s not found", c.Adapter))
	}
	adapterConfig.Name = c.Adapter

	// http client
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid timeout: %v", err))
	}
	httpClient := &http.Client{
		Timeout: timeout,
	}

	// translator
	translator := NewTranslator(adapterConfig,
		WithHTTPClient(httpClient),
	)

	if input, found := strings.CutPrefix(rootFlags.UserInput, "S:"); !found {
		result, err := translator.TranslateWord(input)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	} else {
		result, err := translator.TranslateSentense(input)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
