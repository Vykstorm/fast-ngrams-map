package main

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

const prime = 2654435761

// Implements a map where keys are n-grams. Values may be any type of number value.
type FastNgramsMap[T Number] struct {
	// Unigrams storge
	unigrams    []T
	unigramsSet []bool
	// Bigrams storage
	bigrams    []T
	bigramsSet []bool
	// Trigrams storage
	trigrams map[uint32]T
	// Quadrigrams storage
	quadrigrams map[uint32]T
	// Other n-grams ...
	remaining []map[string]T
}

func NewFastNgramsMap[T Number](maxNgramSize int, n int) *FastNgramsMap[T] {
	m := &FastNgramsMap[T]{}
	m.unigrams = make([]T, 256)
	m.unigramsSet = make([]bool, 256)
	for i := 0; i < 256; i++ {
		m.unigramsSet[i] = false
	}

	m.bigrams = make([]T, 256*256)
	m.bigramsSet = make([]bool, 256*256)
	for i := 0; i < 256*256; i++ {
		m.bigramsSet[i] = false
	}

	m.trigrams = make(map[uint32]T)
	m.quadrigrams = make(map[uint32]T)

	m.remaining = make([]map[string]T, 4*maxNgramSize-4)
	for i := 0; i < len(m.remaining); i++ {
		m.remaining[i] = make(map[string]T)
	}

	return m
}

func (m *FastNgramsMap[T]) Put(ngram string, value T) {
	n := len(ngram)
	switch n {
	case 1:
		m.unigrams[ngram[0]] = value
		m.bigramsSet[ngram[0]] = true
	case 2:
		key := uint16(ngram[0]) | uint16(ngram[1])<<8
		m.bigrams[key] = value
		m.bigramsSet[key] = true
	case 3:
		key := (uint32(ngram[0]) | uint32(ngram[1])<<8 | uint32(ngram[2])<<16) * 2654435761
		m.trigrams[key] = value
	case 4:
		key := (uint32(ngram[0]) | uint32(ngram[1])<<8 | uint32(ngram[2])<<16 | uint32(ngram[3])<<24) * 2654435761
		m.quadrigrams[key] = value
	default:
		m.remaining[n][ngram] = value
	}
}

func (m *FastNgramsMap[T]) Get(ngram string) (T, bool) {
	n := len(ngram)
	switch n {
	case 1:
		var zero T
		key := ngram[0]
		if !m.unigramsSet[key] {
			return zero, false
		}
		value := m.unigrams[key]
		return value, true
	case 2:
		key := uint16(ngram[0]) | uint16(ngram[1])<<8
		var zero T
		if !m.bigramsSet[key] {
			return zero, false
		}
		value := m.bigrams[key]
		return value, true
	case 3:
		key := (uint32(ngram[0]) | uint32(ngram[1])<<8 | uint32(ngram[2])<<16) * prime
		value, ok := m.trigrams[key]
		return value, ok
	case 4:
		key := (uint32(ngram[0]) | uint32(ngram[1])<<8 | uint32(ngram[2])<<16 | uint32(ngram[3])<<24) * prime
		value, ok := m.quadrigrams[key]
		return value, ok
	default:
		value, ok := m.remaining[n][ngram]
		return value, ok
	}

}
