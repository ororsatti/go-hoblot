package search

import (
	"math"
	"sort"
	"strings"

	"github.com/ororsatti/go-searchdex/radix"
)

const (
	editWeight      = 0.5
	relevanceWeight = 0.5
)

// temp till i figure out the tags stuff
type Document struct {
	Id      string
	Content string
}

type docScore struct {
	Key   string
	Score float32
}

type invertedIndex map[string]float32

type termInformation struct {
	docsFreq invertedIndex
}

type SearchIndex struct {
	smap     *radix.SearchableMap
	docCount int
}

func New(docs []Document) *SearchIndex {
	index := SearchIndex{
		smap:     radix.NewSearchableMap(),
		docCount: len(docs),
	}

	for _, doc := range docs {
		index.indexDocument(doc)
	}

	return &index
}

func (index *SearchIndex) indexDocument(doc Document) error {
	termsFreq, err := getTermsFreq(doc.Content)
	if err != nil {
		return err
	}

	for term, freq := range termsFreq {
		termInfo := index.getTermInfo(term)

		if termInfo == nil {
			index.smap.Set(term, &termInformation{
				docsFreq: invertedIndex{
					doc.Id: calculateTf(freq, len(doc.Content)),
				},
			})
		} else {
			termInfo.docsFreq[doc.Id] = calculateTf(freq, len(doc.Content))
			// update the score
		}
	}

	return nil
}

func (index *SearchIndex) Search(query string, maxDistance int) []string {
	queryTerms := strings.Fields(strings.ToLower(query))
	if len(queryTerms) == 0 {
		return []string{}
	}

	relevantDocs := make(map[string]float32)

	for _, queryTerm := range queryTerms {
		fuzzyMatches := index.smap.FuzzyGet(queryTerm, maxDistance)

		var bestMatch *radix.SearchResult
		for _, match := range fuzzyMatches {
			if bestMatch == nil || match.Distance < bestMatch.Distance {
				bestMatch = &match
			}
		}

		if bestMatch == nil {
			continue
		}

		termInfo := bestMatch.Data.(*termInformation)
		if termInfo == nil {
			continue
		}

		idf := calculateIdf(index.docCount, len(termInfo.docsFreq))

		edSimilarity := 1.0 - (float32(bestMatch.Distance) / float32(len(queryTerm)))

		for docKey, tf := range termInfo.docsFreq {
			combinedTermScore := (relevanceWeight * (tf * idf)) + (editWeight * edSimilarity)
			relevantDocs[docKey] += combinedTermScore
		}
	}

	var sortedResults []docScore
	for key, score := range relevantDocs {
		sortedResults = append(sortedResults, docScore{Key: key, Score: score})
	}

	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Score > sortedResults[j].Score
	})

	finalKeys := make([]string, len(sortedResults))
	for i, res := range sortedResults {
		finalKeys[i] = res.Key
	}

	return finalKeys
}

func (index *SearchIndex) IndexDocument(doc Document) error {
	index.docCount++

	if err := index.indexDocument(doc); err != nil {
		index.docCount--
		return err
	}

	return nil
}

func getTermsFreq(content string) (map[string]int, error) {
	tokens := strings.Fields(strings.ToLower(content))

	termsFreq := make(map[string]int)

	for _, tok := range tokens {
		termsFreq[tok]++
	}

	return termsFreq, nil
}

func (index *SearchIndex) getTermInfo(key string) *termInformation {
	data := index.smap.Get(key)
	if data == nil {
		return nil
	}

	termInfoPtr := data.(*termInformation)

	return termInfoPtr
}

func Search(query string) []radix.SearchResult {
	return nil
}

func calculateTf(termFreq, wordCount int) float32 {
	return float32(termFreq) / float32(wordCount)
}

func calculateIdf(docCount, relaventDocCount int) float32 {
	return float32(math.Log(float64(docCount) / float64(relaventDocCount)))
}
