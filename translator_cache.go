package main

import (
	"os"
	"path/filepath"
)

type Cache struct {
	dir string
}

func NewCache(dir string) *Cache {
	return &Cache{
		dir: dir,
	}
}

func (c *Cache) Load(word string) (string, error) {
	filePath := filepath.Join(c.dir, word+".json")

	result, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c *Cache) Save(word, result string) error {
	filePath := filepath.Join(c.dir, word+".json")

	if err := os.MkdirAll(filepath.Dir(c.dir), 0o755); err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(result), 0o644)
}
