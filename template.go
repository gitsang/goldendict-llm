package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"regexp"
)

//go:embed static/sentence.html.tmpl
var sentenceTemplateString string

var sentenceTemplate = template.Must(template.New("sentence.html.tmpl").Parse(sentenceTemplateString))

//go:embed static/user-input.md
var userInputTemplateString string

var userInputTemplate = template.Must(template.New("user-input.md").Parse(userInputTemplateString))

type SentenceEntry struct {
	Sentence    string
	Translation string
}

func RenderSentenceTemplateToString(sentence, translation string) (string, error) {
	var buf bytes.Buffer
	if err := sentenceTemplate.Execute(&buf,
		SentenceEntry{
			Sentence:    sentence,
			Translation: translation,
		}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

//go:embed static/word.html.tmpl
var wordTemplateString string

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}

var wordTemplate = template.Must(template.New("word.html.tmpl").Funcs(funcMap).Parse(wordTemplateString))

type WordEntry struct {
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

func RemoveMarkdownCodeBlockTags(content string) string {
	startRegex := regexp.MustCompile("(?m)^```.*\\n")
	endRegex := regexp.MustCompile("(?m)^```\\s*$")
	content = endRegex.ReplaceAllString(content, "")
	content = startRegex.ReplaceAllString(content, "")
	return content
}

func ProcessWordResponse(content string) (string, error) {
	entry, err := ParseContentToWordEntry(RemoveMarkdownCodeBlockTags(content))
	if err != nil {
		return "", nil
	}
	return RenderWordTemplateToString(entry)
}

func RenderWordTemplateToString(data *WordEntry) (string, error) {
	var buf bytes.Buffer
	if err := wordTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderUserInputTemplateToString(userInput string) (string, error) {
	var buf bytes.Buffer
	if err := userInputTemplate.Execute(&buf, userInput); err != nil {
		return "", err
	}
	return buf.String(), nil
}
