package main

import (
	"os"
	"testing"
)

func TestRenderWordTemplate(t *testing.T) {
	wordEntry := GetSampleWordEntry()

	f, err := os.Create("/tmp/test_render.html")
	if err != nil {
		t.Fatalf("Create file failed: %v", err)
	}
	defer f.Close()

	err = RenderWordTemplate(wordEntry, f)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	t.Log("Render success")
}

func TestRenderWordTemplateToString(t *testing.T) {
	wordEntry := GetSampleWordEntry()

	html, err := RenderWordTemplateToString(wordEntry)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if len(html) == 0 {
		t.Fatal("Rendered html is empty")
	}

	t.Logf("Render success: %d", len(html))
}
