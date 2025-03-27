package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRenderWordTemplateToString(t *testing.T) {
	wordEntry := GetSampleWordEntry()

	wordEntryJsonBytes, _ := json.MarshalIndent(wordEntry, "", "  ")
	fmt.Println(string(wordEntryJsonBytes))

	html, err := RenderWordTemplateToString(wordEntry)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if len(html) == 0 {
		t.Fatal("Rendered html is empty")
	}

	t.Logf("Render success: %d", len(html))
}
