package jsonstore

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotFound = errors.New("key not found")

type Store[T any] struct {
	mu   sync.RWMutex
	path string
	data map[string]json.RawMessage
}

func NewStore[T any](path string) (*Store[T], error) {
	s := &Store[T]{
		path: path,
		data: make(map[string]json.RawMessage),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store[T]) load() error {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(s.path, []byte("{}"), 0o644)
	}
	if err != nil {
		return err
	}
	if len(b) == 0 {
		s.data = make(map[string]json.RawMessage)
		return nil
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store[T]) persist() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store[T]) Set(key string, v T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.data[key] = b
	return s.persist()
}

func (s *Store[T]) Get(key string) (T, error) {
	s.mu.RLock()
	raw, ok := s.data[key]
	s.mu.RUnlock()

	var zero T
	if !ok {
		return zero, ErrNotFound
	}

	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return zero, err
	}
	return v, nil
}

func (s *Store[T]) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.persist()
}

func (s *Store[T]) Keys() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys, nil
}
