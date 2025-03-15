package src

import (
	"math/rand"
	"testing"
)

// TestMapInt tests the map implementation with T as int.
func TestDefaultMapInt(t *testing.T) {
	CheckMap[int](t, NewDefaultMap[int](0), func() int {
		return rand.Intn(1000) // Generate a random int value.
	})
}

// TestMapFloat64 tests the map implementation with T as float64.
func TestDefaultMapFloat64(t *testing.T) {
	CheckMap[float64](t, NewDefaultMap[float64](0), func() float64 {
		return rand.Float64() * 1000 // Generate a random float64 value.
	})
}

// TestMapUint32 tests the map implementation with T as uint32.
func TestDefaultMapUint32(t *testing.T) {
	CheckMap[uint32](t, NewDefaultMap[uint32](0), func() uint32 {
		return uint32(rand.Intn(1000)) // Generate a random uint32 value.
	})
}
func TestDefaultMapConcurrentReads(t *testing.T) {
	CheckMapConcurrentReads(t, NewDefaultMap[int](0),
		func() int { return rand.Intn(1000) })
}
