package main

import (
	_ "embed"
	"html/template"
)

//go:embed static/user-input.md
var userInputTemplateString string

var userInputTemplate = template.Must(template.New("user-input.md").Parse(userInputTemplateString))

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}
