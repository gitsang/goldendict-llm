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

func RenderSentenceTemplateToString(sentence, translation, adapterName, model, duration string) (string, error) {
	var buf bytes.Buffer
	if err := sentenceTemplate.Execute(&buf,
		SentenceEntry{
			AdapterName: adapterName,
			Model:       model,
			Duration:    duration,
			Sentence:    sentence,
			Translation: translation,
		}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
