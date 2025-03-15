package src

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testMapGeneric runs a suite of tests on a Map[T] using a random generator
// function 'randVal' to produce values of type T.
func CheckMap[T any](t *testing.T, m Map[T], randVal func() T) {

	// Extended candidate runes for keys covering different Unicode ranges.
	candidateRunes := []rune{
		'a', 'b', 'c', // Basic Latin letters.
		'中', '文', // Chinese characters.
		'α', 'β', 'γ', // Greek letters.
		'😀', '😃', '😄', // Emojis.
		'न', 'प', // Devanagari.
		'é', 'ñ', // Latin letters with accents.
		'𐍈', // Gothic letter.
	}
	totalKeys := 0

	// For each candidate rune, generate keys of length 1 to 7.
	for _, r := range candidateRunes {
		for length := 1; length <= 7; length++ {
			key := strings.Repeat(string(r), length)

			// Check that the key is not already present.
			assert.False(t, m.Has(key), "Expected key %q to not exist initially", key)
			_, ok := m.Get(key)
			assert.False(t, ok, "Expected Get for key %q to return false initially", key)

			// Generate a random value for the entry.
			val := randVal()
			m.Put(key, val)
			totalKeys++

			// Validate that the key now exists with the expected value.
			assert.True(t, m.Has(key), "Expected key %q to exist after insertion", key)
			gotVal, ok := m.Get(key)
			assert.True(t, ok, "Expected Get for key %q to return true after insertion", key)
			assert.Equal(t, val, gotVal, "Expected value for key %q to be %v", key, val)

			assert.Equal(t, totalKeys, m.Len(), "Expected map length to equal total inserted keys")
		}
	}

	// Validate that the length of the map equals the total keys inserted.
	assert.Equal(t, totalKeys, m.Len(), "Expected map length to equal total inserted keys")

	// Test updating an existing key.
	testKey := strings.Repeat(string(candidateRunes[0]), 3) // e.g. "aaa"
	newValue := randVal()
	m.Put(testKey, newValue)
	// Length should remain unchanged.
	assert.Equal(t, totalKeys, m.Len(), "Expected map length to remain unchanged after update")
	gotValue, ok := m.Get(testKey)
	assert.True(t, ok, "Expected updated key %q to be present", testKey)
	assert.Equal(t, newValue, gotValue, "Expected updated value for key %q to be %v", testKey, newValue)

	// Verify that querying a non-existent key returns false.
	nonExistentKey := "nonexistent"
	_, ok = m.Get(nonExistentKey)
	assert.False(t, ok, "Expected Get for non-existent key %q to return false", nonExistentKey)
}

func CheckMapConcurrentReads[T any](t *testing.T, m Map[T], randVal func() T) {
	// Populate the map with a fixed number of entries.
	const numEntries = 100
	values := make([]T, numEntries)

	for i := 0; i < numEntries; i++ {
		key := fmt.Sprintf("key%d", i)
		values[i] = randVal()
		m.Put(key, values[i])
	}

	// Define the number of concurrent reader goroutines and iterations per goroutine.
	const numReaders = 20
	const numIterations = 1000

	var wg sync.WaitGroup

	// Launch concurrent goroutines that perform Len() and Get() operations.
	for readerID := 0; readerID < numReaders; readerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Call Len() concurrently.
				length := m.Len()
				assert.Equal(t, numEntries, length, "Goroutine %d: expected length %d, got %d", id, numEntries, length)

				// Call Get() concurrently on a key selected in a round-robin fashion.
				key := fmt.Sprintf("key%d", j%numEntries)
				val, ok := m.Get(key)
				assert.True(t, ok, "Goroutine %d: expected key %q to exist", id, key)
				assert.Equal(t, values[j%numEntries], val, "Goroutine %d: expected value %d for key %q, got %v", id, j%numEntries, key, val)
			}
		}(readerID)
	}

	// Wait for all goroutines to complete.
	wg.Wait()
}
