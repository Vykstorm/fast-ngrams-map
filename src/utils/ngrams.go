package utils

import (
	"iter"
	"slices"
	"unicode/utf8"
)

/*
*
Calculate the number of n-grams with the specified lengths in the given text.
*/
func GetNgramsCount(text string, gramSizes []int) int {
	n := utf8.RuneCountInString(text)
	m := 0
	for _, i := range gramSizes {
		if i > n {
			continue
		}
		m += n - i + 1
	}
	return m
}

/*
*
Iterates over the all the n-grams (shorter n-grams are iterated first) with the specified lengths in the given text.
*/
func IterNgrams(text string, gramSizes []int) iter.Seq[string] {
	gramSizesSorted := make([]int, len(gramSizes))
	copy(gramSizesSorted, gramSizes)
	slices.Sort(gramSizesSorted)

	chars := []rune(text)
	n := len(chars)

	return func(yield func(string) bool) {
		for _, ngramSize := range gramSizes {
			for j := 0; j < n-ngramSize+1; j++ {
				ngram := string(chars[j : j+ngramSize])
				if !yield(ngram) {
					return
				}
			}
		}
	}
}

/*
Get all the n-grams with the specified lengths of the given text as an array.
(shorter n-grams are placed first)
*/
func GetNgrams(text string, gramSizes []int) []string {
	numNgrams := GetNgramsCount(text, gramSizes)
	ngrams := make([]string, numNgrams)
	i := 0
	for ngram := range IterNgrams(text, gramSizes) {
		ngrams[i] = ngram
		i++
	}
	return ngrams
}

/*
*
Get a language profile. A map that contains the n-grams found in the given file
with their occurrence frequencies in a map.
*/
func GetLanguageProfile(filePath string, gramSizes []int) (map[string]int, error) {
	samples, err := ReadLines(filePath)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, sample := range samples {
		for ngram := range IterNgrams(sample, gramSizes) {
			if _, ok := counts[ngram]; !ok {
				counts[ngram] = 1
			} else {
				counts[ngram] += 1
			}
		}
	}
	return counts, nil
}

/*
Given the ngrams and their frequencies, returns the n most frequent n-grams
*/
func GetNgramsRanking(profile map[string]int, n int) []string {
	ranking := make([]string, len(profile))
	i := 0
	for ngram := range profile {
		ranking[i] = ngram
		i += 1
	}
	slices.SortFunc(ranking, func(a string, b string) int {
		return profile[a] - profile[b]
	})
	slices.Reverse(ranking)
	return ranking
}
