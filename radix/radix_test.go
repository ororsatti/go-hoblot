package radix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	testCases := []struct {
		keys     []string
		expected []string
	}{
		{
			keys:     []string{"bird"},
			expected: []string{"bird"},
		},
		{
			keys:     []string{"bird", "barb"},
			expected: []string{"b", "ird", "arb"},
		},
		{
			keys:     []string{"bird", "birds"},
			expected: []string{"bird", "s"},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case %d", i+1), func(tt *testing.T) {
			smap := NewSearchableMap()

			for _, key := range testCase.keys {
				smap.Set(key, "")
			}

			for _, key := range testCase.expected {
				assert.NotNil(tt, smap.root.lookup(key))
			}
		})
	}
}

func TestGet(t *testing.T) {
	str := "the quick brown fox jumped over the red fence"
	keys := strings.Fields(str)
	smap := NewSearchableMap()

	for _, key := range keys {
		smap.Set(key, true)
	}

	for _, key := range keys {
		assert.NotNil(t, smap.Get(key))
	}
}

func TestGetNotExist(t *testing.T) {
	smap := NewSearchableMap()
	smap.Set("bird", true)

	assert.Nil(t, smap.Get("what"))
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name        string
		keys        []string
		keyToDelete string
	}{
		{
			name:        "delete exist",
			keys:        []string{"bird", "brew", "brand"},
			keyToDelete: "bird",
		},
		{
			name:        "delete all",
			keys:        []string{"bird"},
			keyToDelete: "bird",
		},
		{
			name:        "delete not exist",
			keys:        []string{"bird", "brew", "brand"},
			keyToDelete: "band",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			smap := NewSearchableMap()
			for _, key := range testCase.keys {
				smap.Set(key, true)
			}

			smap.Delete(testCase.keyToDelete)

			for _, key := range testCase.keys {
				if key != testCase.keyToDelete {
					assert.NotNil(tt, smap.Get(key))
				} else {
					assert.Nil(tt, smap.Get(key))
				}
			}
		})
	}
}

func TestFuzzyGet(t *testing.T) {
	terms := []string{"wonder", "ponder", "wondering", "ball", "inter"}
	testCases := []struct {
		name        string
		key         string
		maxDistance int
		results     map[string]SearchResult
	}{
		{
			name:        "no results",
			key:         "wombat",
			maxDistance: 1,
			results:     map[string]SearchResult{},
		},
		{
			name:        "multiple results",
			key:         "winter",
			maxDistance: 3,
			results: map[string]SearchResult{
				"inter": {
					Distance: 1,
					Data:     true,
				},
				"wonder": {
					Distance: 2,
					Data:     true,
				},
				"ponder": {
					Distance: 3,
					Data:     true,
				},
			},
		},
		{
			name:        "single result",
			key:         "hall",
			maxDistance: 1,
			results: map[string]SearchResult{
				"ball": {
					Distance: 1,
					Data:     true,
				},
			},
		},
	}

	smap := NewSearchableMap()

	for _, term := range terms {
		smap.Set(term, true)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			assert.Equal(tt, testCase.results, smap.FuzzyGet(testCase.key, testCase.maxDistance))
		})
	}
}
