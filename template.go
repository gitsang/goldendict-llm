package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"regexp"
)

//go:embed static/user-input.md
var userInputTemplateString string

var userInputTemplate = template.Must(template.New("user-input.md").Parse(userInputTemplateString))

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}

func RenderUserInputTemplateToString(userInput string) (string, error) {
	var buf bytes.Buffer
	if err := userInputTemplate.Execute(&buf, userInput); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RemoveMarkdownCodeBlockTags(content string) string {
	startRegex := regexp.MustCompile("(?m)^```.*\\n")
	endRegex := regexp.MustCompile("(?m)^```\\s*$")
	content = endRegex.ReplaceAllString(content, "")
	content = startRegex.ReplaceAllString(content, "")
	return content
}
