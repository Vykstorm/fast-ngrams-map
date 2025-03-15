package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/Vykstorm/fast-ngrams-map/src"
	"github.com/Vykstorm/fast-ngrams-map/src/utils"
)

/**
The benchmark evaluates the time performance of different map implementations used to store
to store n-grams and their frequencies of occurrence in arbitrary language corpora:

- For each language, obtain the K most frequent n-grams in the map.
- Evaluate the performance of the map Get() operation by calling it for each extracted n-gram.
from a set of text samples written in the same language.
*/

const MaxBenchmarkTime = time.Duration(time.Second) * 20

const NgramsDefaultMap = 1
const NgramsFastMap = 2

var MapImplementationName = map[int]string{
	NgramsDefaultMap: "Golang built-in map",
	NgramsFastMap:    "N-grams fast map",
}

func createMap(mapImplementation int, ngramsRanking []string, profile map[string]int) src.Map[int] {
	var m src.Map[int]
	switch mapImplementation {
	case NgramsFastMap:
		m = src.NewFastNgramsMap[int](len(ngramsRanking))
	default:
		m = src.NewDefaultMap[int](len(ngramsRanking))
	}
	for _, ngram := range ngramsRanking {
		m.Put(ngram, profile[ngram])
	}
	return m
}

func getCorporaFilePathForLanguage(language string) string {
	return filepath.Join("corpora", language+".txt")
}

func RunBenchmark(language string, mapImplementation int, gramSizes []int, rankingSize int, numSamplesToTest int) {
	// Build language profile
	profile, err := utils.GetLanguageProfile(getCorporaFilePathForLanguage(language), gramSizes)
	if err != nil {
		log.Fatal(err)
	}
	// Get ngrams ranking
	ngramsRanking := utils.GetNgramsRanking(profile, rankingSize)

	// Build map
	m := createMap(mapImplementation, ngramsRanking, profile)

	// Get language random samples
	samples, err := utils.ReadNRandomLines(getCorporaFilePathForLanguage(language), numSamplesToTest)
	if err != nil {
		log.Fatal(err)
	}

	// Get all n-grams
	ngrams := make([]string, 0)
	for _, sample := range samples {
		ngrams = append(ngrams, utils.GetNgrams(sample, gramSizes)...)
	}

	var elapsedTime float64 = 0
	totalNgramsProcessed := 0

	for {
		t := time.Now()
		for _, ngram := range ngrams {
			m.Get(ngram)
		}
		diff := time.Since(t)
		totalNgramsProcessed += len(ngrams)
		elapsedTime += diff.Seconds()
		if elapsedTime >= MaxBenchmarkTime.Seconds() {
			break
		}
	}

	elapsedTimePerNgramProcessed := elapsedTime / float64(totalNgramsProcessed)

	var benchmarkResultInfo = struct {
		MapImplementation        string `json:"implementation"`
		Language                 string `json:"language"`
		GramSizes                []int  `json:"gramSizes"`
		RankingSize              int    `json:"ngramsRankingSize"`
		TimePerLookupNanoseconds int64  `json:"timePerLookupNs"`
		TimePerLookupStr         string `json:"timePerLookupStr"`
	}{
		MapImplementation:        MapImplementationName[mapImplementation],
		Language:                 language,
		GramSizes:                gramSizes,
		RankingSize:              rankingSize,
		TimePerLookupNanoseconds: int64(elapsedTimePerNgramProcessed * 1e9),
		TimePerLookupStr:         time.Duration(int64(elapsedTimePerNgramProcessed * 1e9)).String(),
	}
	info, err := json.MarshalIndent(benchmarkResultInfo, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(info))

}

func main() {
	RunBenchmark("spa-latn", NgramsDefaultMap, []int{1, 2, 3}, 30000, 50)
	RunBenchmark("spa-latn", NgramsFastMap, []int{1, 2, 3}, 30000, 50)
}
