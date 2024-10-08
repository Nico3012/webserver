package jsoncache

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type JsonCache string

func (jsonCache JsonCache) loadMap() (map[string]string, error) {
	m := make(map[string]string)

	_, err := os.Stat(string(jsonCache))
	if os.IsNotExist(err) {
		return m, nil
	}

	file, err := os.ReadFile(string(jsonCache))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (jsonCache JsonCache) saveMap(m map[string]string) error {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(string(jsonCache), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// exported functions:

func (jsonCache JsonCache) Get(ctx context.Context, key string) ([]byte, error) {
	m, err := jsonCache.loadMap()
	if err != nil {
		return []byte{}, err
	}

	value, exists := m[key]
	if !exists {
		return []byte{}, autocert.ErrCacheMiss
	}

	return []byte(value), nil
}

func (jsonCache JsonCache) Put(ctx context.Context, key string, data []byte) error {
	m, err := jsonCache.loadMap()
	if err != nil {
		return err
	}

	m[key] = string(data)

	err = jsonCache.saveMap(m)
	if err != nil {
		return err
	}

	return nil
}

func (jsonCache JsonCache) Delete(ctx context.Context, key string) error {
	m, err := jsonCache.loadMap()
	if err != nil {
		return err
	}

	_, exists := m[key]
	if !exists {
		return nil
	}

	delete(m, key)

	err = jsonCache.saveMap(m)
	if err != nil {
		return err
	}

	return nil
}
