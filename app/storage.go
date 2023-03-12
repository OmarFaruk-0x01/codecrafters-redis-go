package main

import (
	"errors"
	"time"
)

var (
	errNoKeyFound = errors.New("invalid key")
)

type Data struct {
	value string
	exp   time.Time
}

type Storage struct {
	data map[string]Data
}

func NewStorage() *Storage {
	return &Storage{
		data: map[string]Data{},
	}
}

func (s *Storage) SetItem(key string, data *Data) {
	s.data[key] = *data
}

func (s *Storage) GetItem(key string) (Data, error) {
	if data, ok := s.data[key]; ok {
		return data, nil
	}
	return Data{}, errNoKeyFound
}
