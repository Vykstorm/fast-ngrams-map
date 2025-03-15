package src

/**
Implementation of Map interface with the following properties:

- Entries with keys having one or two bytes in length are stored in arrays of size 256 and 65536 respectively.
  Direct hash mapping is used to boost lookup time.

- Entries with keys of 3-bytes and 4-bytes in length are stored in separate hash maps.
  The underline hash maps use 32-bit integers as keys.

  The original key (string of 3 or 4 bytes) is converted to a 32-bit integer value
  and multiplied by a prime number to increase hash sparsity.

- Entries with keys with size of 5 or more bytes are stored in a regular hash map.
*/

const prime = 2654435761

type FastNgramsMap[T any] struct {
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
	remaining map[string]T
}

func NewFastNgramsMap[T any](size int) *FastNgramsMap[T] {
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

	m.remaining = make(map[string]T, size)

	return m
}

func (m *FastNgramsMap[T]) Put(ngram string, value T) {
	n := len(ngram)
	switch n {
	case 1:
		m.unigrams[ngram[0]] = value
		m.unigramsSet[ngram[0]] = true
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
		m.remaining[ngram] = value
	}
}

func (m *FastNgramsMap[T]) Has(ngram string) bool {
	_, ok := m.Get(ngram)
	return ok
}

func (m *FastNgramsMap[T]) Len() int {
	count := 0
	for _, isSet := range m.unigramsSet {
		if isSet {
			count += 1
		}
	}
	for _, isSet := range m.bigramsSet {
		if isSet {
			count += 1
		}
	}
	count += len(m.trigrams)
	count += len(m.quadrigrams)
	count += len(m.remaining)

	return count
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
		value, ok := m.remaining[ngram]
		return value, ok
	}

}
