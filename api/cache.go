package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SaveToCache func(filePath string, data []byte) error

// saveToCache creates the cache directory and file, then saves the bytes to
// the cache file.
func saveToCache(filePath string, data []byte) error {
	// Create cache directory
	if err := os.MkdirAll("cache", os.ModePerm); err != nil {
		return fmt.Errorf("error in saveToCache while creating cache directory: %w", err)
	}

	// Create cache file
	f, fErr := os.Create(filePath)
	if fErr != nil {
		return fmt.Errorf("error in saveToCache while creating cache file: %w", fErr)
	}

	// Write to cache file
	_, writeToFileErr := f.Write(data)
	if writeToFileErr != nil {
		f.Close()
		return fmt.Errorf("error in saveToCache writing to cache file: %w", writeToFileErr)
	}

	// Close cache file
	if err := f.Close(); err != nil {
		return fmt.Errorf("error in saveToCache while closing cache file: %w", err)
	}

	fmt.Println("File written successfully.")

	return nil
}

// readCacheAs reads a file's raw bytes and unmarshals them into the provided
// type.
func readCacheAs[T any](filePath string) (T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("file doesn't exist or is unreadable: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("could not decode cached JSON: %w", err)
	}

	return result, nil
}

type ReadRssFromCache func(filePath string) (RssResponse, error)

// readRssFromCache reads and decodes cached RSS data from the specified file path.
func readRssFromCache(filePath string) (RssResponse, error) {
	return readCacheAs[RssResponse](filePath)
}
