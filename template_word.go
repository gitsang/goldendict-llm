package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
)

//go:embed static/word-prompt.md
var WordPrompt string

//go:embed static/word.html.tmpl
var wordTemplateString string

var wordTemplate = template.Must(template.New("word.html.tmpl").Funcs(funcMap).Parse(wordTemplateString))

type WordEntry struct {
	AdapterName string
	Model       string
	Duration    string

	Word                  string
	Pronunciation         Pronunciation
	PartOfSpeech          string
	GrammaticalInfo       string
	Definitions           []Definition
	Idioms                []Idiom
	RelatedWords          *RelatedWord
	AlternativeDefinition *AlternativeDefinition
}

type Pronunciation struct {
	Syllables   string
	Phonetic    string
	Alternative string
}

type Definition struct {
	English  string
	Chinese  string
	Synonym  string
	Examples []Example
}

type Example struct {
	English string
	Chinese string
}

type Idiom struct {
	Phrase   string
	Style    string
	English  string
	Chinese  string
	Examples []Example
}

type RelatedWord struct {
	Word            string
	PartOfSpeech    string
	GrammaticalInfo string
}

type AlternativeDefinition struct {
	Word            string
	PartOfSpeech    string
	GrammaticalInfo string
	Definitions     []Definition
}

func ParseContentToWordEntry(content string) (*WordEntry, error) {
	wordEntry := WordEntry{}
	err := json.Unmarshal([]byte(content), &wordEntry)
	if err != nil {
		return nil, err
	}
	return &wordEntry, nil
}

func RenderWordTemplateToString(entry *WordEntry) (string, error) {
	var buf bytes.Buffer
	if err := wordTemplate.Execute(&buf, entry); err != nil {
		return "", err
	}
	return buf.String(), nil
}
