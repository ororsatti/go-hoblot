package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type doc struct {
	Id      string `index:"id"`
	Content string `index:"text"`
}

func TestIndexDocument(t *testing.T) {
	documents := []doc{
		{
			Id:      "doc1",
			Content: "a b",
		},
		{
			Id:      "doc2",
			Content: "a b",
		},
		{
			Id:      "doc3",
			Content: "a c c",
		},
	}

	testCases := []struct {
		testName  string
		key       string
		docsCount int
	}{
		{
			testName:  "key a",
			key:       "a",
			docsCount: 3,
		},
		{
			testName:  "key b",
			key:       "b",
			docsCount: 2,
		},
		{
			testName:  "key c",
			key:       "c",
			docsCount: 1,
		},
	}

	anys := make([]any, 3)
	for i := range len(anys) {
		anys[i] = documents[i]
	}

	index := New(anys)
	assert.Equal(t, index.docCount, 3)
	for _, testCase := range testCases {
		t.Run(testCase.testName, func(tt *testing.T) {
			assert.NotNil(t, index.smap.Get(testCase.key))

			termInfo := index.getTermInfo(testCase.key)

			assert.Equal(t, len(termInfo.docsFreq), testCase.docsCount)
		})
	}

	index.IndexDocument(doc{
		Id:      "doc4",
		Content: "d",
	})

	assert.Equal(t, index.docCount, 4)
	termInfo := index.getTermInfo("d")
	assert.NotNil(t, termInfo)
	assert.Equal(t, len(termInfo.docsFreq), 1)
}

func TestSearch(t *testing.T) {
	documents := []doc{
		{
			Id:      "doc1",
			Content: "the quick brown fox",
		},
		{
			Id:      "doc2",
			Content: "jumps over the lazy dog",
		},
		{
			Id:      "doc3",
			Content: "a quick brown dog",
		},
	}

	anys := make([]any, 3)
	for i := range len(anys) {
		anys[i] = documents[i]
	}
	index := New(anys)

	testCases := []struct {
		name        string
		query       string
		maxDistance int
		expected    []string
	}{
		{
			name:        "exact match",
			query:       "fox",
			maxDistance: 0,
			expected:    []string{"doc1"},
		},
		{
			name:        "fuzzy match",
			query:       "quik",
			maxDistance: 1,
			expected:    []string{"doc1", "doc3"},
		},
		{
			name:        "multiple terms",
			query:       "brown dog",
			maxDistance: 0,
			expected:    []string{"doc3", "doc2", "doc1"},
		},
		{
			name:        "no match",
			query:       "cat",
			maxDistance: 0,
			expected:    []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := index.Search(tc.query, tc.maxDistance)
			assert.ElementsMatch(t, tc.expected, results)
		})
	}
}
