package src

/**
Default implementation of the Map interface with built-in Golang map.
*/

type DefaultMap[T any] struct {
	entries map[string]T
}

func (m *DefaultMap[T]) Get(key string) (T, bool) {
	value, ok := m.entries[key]
	return value, ok
}

func (m *DefaultMap[T]) Has(key string) bool {
	_, ok := m.entries[key]
	return ok
}

func (m *DefaultMap[T]) Put(key string, value T) {
	m.entries[key] = value
}

func (m *DefaultMap[T]) Len() int {
	return len(m.entries)
}

func NewDefaultMap[T any](size int) *DefaultMap[T] {
	return &DefaultMap[T]{
		entries: make(map[string]T, size),
	}
}
