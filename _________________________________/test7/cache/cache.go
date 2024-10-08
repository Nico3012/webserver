package cache

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type jsonMap map[string]string
type Cache string

func (cache Cache) Get(ctx context.Context, key string) ([]byte, error) {
	jsonMap, err := loadJsonMap(string(cache))
	if err != nil {
		return []byte{}, err
	}

	value, exists := jsonMap[key]
	if !exists {
		return []byte{}, autocert.ErrCacheMiss
	}

	return []byte(value), nil
}

func (cache Cache) Put(ctx context.Context, key string, data []byte) error {
	jsonMap, err := loadJsonMap(string(cache))
	if err != nil {
		return err
	}

	jsonMap[key] = string(data)

	err = saveJsonMap(string(cache), jsonMap)
	if err != nil {
		return err
	}

	return nil
}

func (cache Cache) Delete(ctx context.Context, key string) error {
	jsonMap, err := loadJsonMap(string(cache))
	if err != nil {
		return err
	}

	_, exists := jsonMap[key]
	if !exists {
		// return nil if key does not exist
		return nil
	}

	delete(jsonMap, key)

	err = saveJsonMap(string(cache), jsonMap)
	if err != nil {
		return err
	}

	return nil
}

// helper functions

func loadJsonMap(cache string) (jsonMap, error) {
	jsonMap := make(jsonMap)

	_, err := os.Stat(cache)
	if os.IsNotExist(err) {
		// return empty jsonMap if file not exists
		return jsonMap, nil
	}

	file, err := os.ReadFile(cache)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func saveJsonMap(cache string, jsonMap jsonMap) error {
	data, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cache, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
