package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexDocument(t *testing.T) {
	documents := []Document{
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

	index := New(documents)
	assert.Equal(t, index.docCount, 3)
	for _, testCase := range testCases {
		t.Run(testCase.testName, func(tt *testing.T) {
			assert.NotNil(t, index.smap.Get(testCase.key))

			termInfo := index.getTermInfo(testCase.key)

			assert.Equal(t, len(termInfo.docsFreq), testCase.docsCount)
		})
	}
}
