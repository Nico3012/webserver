package main

import (
	"encoding/json"
	"errors"
	"os"
)

var ErrKeyNotFound = errors.New("key not found")

// Cache type to hold key-value pairs
type Cache map[string]string

// Load cache from JSON file
func loadCache(cacheFile string) (Cache, error) {
	cache := make(Cache)

	_, err := os.Stat(cacheFile)

	// Check if the file exists
	if os.IsNotExist(err) {
		// return an empty cache if the file doesn't exist
		return cache, nil
	}

	// Read the file content
	fileContent, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the content into the cache map
	err = json.Unmarshal(fileContent, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Save cache to JSON file
func saveCache(cacheFile string, cache Cache) error {
	// Marshal cache into JSON
	data, err := json.MarshalIndent(cache, "", "    ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Put(cacheFile string, key string, value string) error {
	cache, err := loadCache(cacheFile)
	if err != nil {
		return err
	}

	// Add or update the key-value pair
	cache[key] = value

	// Save the updated cache to file
	err = saveCache(cacheFile, cache)
	if err != nil {
		return err
	}

	return nil
}

func Get(cacheFile string, key string) (string, error) {
	cache, err := loadCache(cacheFile)
	if err != nil {
		return "", err
	}

	// Check if the key exists
	value, exists := cache[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func Delete(cacheFile string, key string) error {
	cache, err := loadCache(cacheFile)
	if err != nil {
		return err
	}

	// Check if the key exists
	_, exists := cache[key]
	if !exists {
		return ErrKeyNotFound
	}

	// Delete the key
	delete(cache, key)

	// Save the updated cache to file
	err = saveCache(cacheFile, cache)
	if err != nil {
		return err
	}

	return nil
}
