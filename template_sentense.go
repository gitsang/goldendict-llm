package main

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed static/sentence-prompt.md
var SentencePrompt string

//go:embed static/sentence.html.tmpl
var sentenceTemplateString string

var sentenceTemplate = template.Must(template.New("sentence.html.tmpl").Parse(sentenceTemplateString))

type SentenceEntry struct {
	AdapterName string
	Model       string
	Duration    string

	Sentence    string
	Translation string
}

func RenderSentenceTemplateToString(entry *SentenceEntry) (string, error) {
	var buf bytes.Buffer
	if err := sentenceTemplate.Execute(&buf, entry); err != nil {
		return "", err
	}
	return buf.String(), nil
}
