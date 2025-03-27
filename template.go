package main

import (
	"bytes"
	"html/template"
	"io"
	"path/filepath"
)

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

func RenderWordTemplate(data *WordEntry, writer io.Writer) error {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl, err := template.New("word.html.tmpl").
		Funcs(funcMap).
		ParseFiles(filepath.Join("static", "word.html.tmpl"))
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func RenderWordTemplateToString(data *WordEntry) (string, error) {
	var buf bytes.Buffer
	if err := RenderWordTemplate(data, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetSampleWordEntry() *WordEntry {
	return &WordEntry{
		Word: "abandon",
		Pronunciation: Pronunciation{
			Syllables:   "a∙ban∙don¹",
			Phonetic:    "E5bAndEn",
			Alternative: "E5bAndEn",
		},
		PartOfSpeech:    "verb",
		GrammaticalInfo: "[T]",
		Definitions: []Definition{
			{
				English: "to leave someone, especially someone you are responsible for",
				Chinese: "抛弃，遗弃〔某人〕",
				Examples: []Example{
					{
						English: "How could she abandon her own child?",
						Chinese: "她怎么能抛弃自己的孩子？",
					},
				},
			},
			{
				English: "to go away from a place, vehicle etc permanently, especially because the situation makes it impossible for you to stay",
				Chinese: "离弃，逃离〔某地方或交通工具等〕",
				Synonym: "leave",
				Examples: []Example{
					{
						English: "We had to abandon the car and walk the rest of the way.",
						Chinese: "我们只好弃车，走完剩下的路。",
					},
					{
						English: "Fearing further attacks, most of the population had abandoned the city.",
						Chinese: "由于害怕遭受更多的袭击，大多数市民已逃离该市。",
					},
				},
			},
			{
				English: "to stop doing something because there are too many problems and it is impossible to continue",
				Chinese: "放弃，中止",
				Examples: []Example{
					{
						English: "The game had to be abandoned due to bad weather.",
						Chinese: "由于天气不好，比赛不得不中止。",
					},
					{
						English: "They abandoned their attempt to recapture the castle.",
						Chinese: "他们放弃了夺回城堡的努力。",
					},
					{
						English: "Because of the fog they abandoned their idea of driving.",
						Chinese: "因为有雾，他们打消了开车去的念头。",
					},
				},
			},
			{
				English: "to stop having a particular idea, belief, or attitude",
				Chinese: "放弃〔信念、信仰或看法〕",
				Examples: []Example{
					{
						English: "They were accused of abandoning their socialist principles.",
						Chinese: "他们被指责放弃了社会主义原则。",
					},
					{
						English: "Rescuers had abandoned all hope of finding any more survivors.",
						Chinese: "营救人员对找到更多生还者已不抱任何希望。",
					},
				},
			},
		},
		Idioms: []Idiom{
			{
				Phrase:  "abandon yourself to sth",
				Style:   "[literary]",
				English: "to feel an emotion so strongly that you let it control you completely",
				Chinese: "沉湎于; 放纵〔感情〕",
				Examples: []Example{
					{
						English: "She abandoned herself to grief.",
						Chinese: "她陷入悲痛之中，不能自拔。",
					},
				},
			},
			{
				Phrase:  "abandon ship",
				English: "to leave a ship because it is sinking",
				Chinese: "〔由于船在下沉而〕弃船（逃生）",
			},
		},
		RelatedWords: &RelatedWord{
			Word:            "abandonment",
			PartOfSpeech:    "noun",
			GrammaticalInfo: "[U]",
		},
		AlternativeDefinition: &AlternativeDefinition{
			Word:            "abandon²",
			PartOfSpeech:    "noun",
			GrammaticalInfo: "[U]",
			Definitions: []Definition{
				{
					English: "if someone does something with abandon, they behave in a careless or uncontrolled way, without thinking or caring about what they are doing",
					Chinese: "尽情; 放任",
					Examples: []Example{
						{
							English: "They drank and smoked with reckless abandon.",
							Chinese: "他们纵情地喝酒抽烟，毫无顾忌。",
						},
					},
				},
			},
		},
	}
}
