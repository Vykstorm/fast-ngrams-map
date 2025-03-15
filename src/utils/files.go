package utils

import (
	"bufio"
	"math/rand/v2"
	"os"
)

// ReadLines reads a text file and returns an array of strings (one per line)
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Rea at most N random lines from the given file.
func ReadNRandomLines(filename string, n int) ([]string, error) {
	samples, err := ReadLines(filename)
	if err != nil {
		return nil, err
	}
	if len(samples) <= n {
		return samples, nil
	}
	indices := rand.Perm(n)
	perm := make([]string, n)
	for i, index := range indices {
		perm[i] = samples[index]
	}
	return perm, nil
}
