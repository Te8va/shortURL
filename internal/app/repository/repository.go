package repository

import (
	"encoding/json"
	"os"
)

type MapStore struct {
	data map[string]string
	file string
}

func NewMapStore(filePath string) (*MapStore, error) {
	store := &MapStore{
		data: make(map[string]string),
		file: filePath,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        file, err := os.Create(filePath)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        data := "{}"
        _, err = file.WriteString(data)
        if err != nil {
            return nil, err
        }
    }

	if err := store.loadFromFile(); err != nil {
		return nil,err
	}

	return store, nil
}

func (s *MapStore) Save(id, url string) error {
	s.data[id] = url

	return s.saveToFile()
}

func (s *MapStore) Get(id string) (string, bool) {
	url, exists := s.data[id]
	return url, exists
}

func (s *MapStore) saveToFile() error {
	file, err := os.OpenFile(s.file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (s *MapStore) loadFromFile() error {
	file, err := os.ReadFile(s.file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, &s.data); err != nil {
		return err
	}
	return nil
}
