package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	errNoKeyFound = errors.New("invalid key")
	errExpireKey  = errors.New("expire key")
)

type Data struct {
	value     string
	createdAt int64
	exp       int64
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
	data.createdAt = time.Now().UnixMilli()
	if data.exp != 0 {
		data.exp = data.createdAt + data.exp
	}
	s.data[key] = *data
}

func (s *Storage) GetItem(key string) (Data, error) {
	data, ok := s.data[key]
	if !ok {
		return Data{}, errNoKeyFound
	}
	fmt.Println(data)
	if data.exp == 0 {
		return data, nil
	}
	currentMil := time.Now().UnixMilli()
	if currentMil > data.exp {
		return Data{}, errExpireKey
	}

	return data, nil

}
