package radix

import (
	"fmt"
	"iter"
	"maps"
)

type SearchableMap struct {
	root *node
}

func NewSearchableMap() *SearchableMap {
	return &SearchableMap{
		root: newNode(),
	}
}

func (s *SearchableMap) Set(key string, data any) *node {
	node := s.root.createPath(key)
	node.Data = data
	return node
}

func (s *SearchableMap) Get(key string) any {
	node := s.root.lookup(key)
	if node.Data != nil {
		return node.Data
	}

	return nil
}

func (s *SearchableMap) Delete(key string) {
	node, path := trackDown(s.root, key, []nodePath{})
	if node == nil || node.Data == nil {
		return
	}

	next, stop := iter.Pull2(maps.All(node.Children))
	defer stop()

	firstKey, firstChild, ok := next()
	if !ok {
		cleanup(path)
	} else {
		merge(path, firstKey, firstChild)
	}
}

func (s *SearchableMap) FuzzyGet(query string, maxDistance int) map[string]SearchResult {
	w := len(query) + 1
	h := w + maxDistance
	matrix := make([]int, w*h)

	for i := range w {
		matrix[i] = i
	}

	for i := range h {
		matrix[i*w] = i
	}

	results := make(map[string]SearchResult)
	s.root.fuzzyRecurse(query, "", matrix, maxDistance, 1, w, results)

	return results
}

func (s *SearchableMap) print() {
	fmt.Println("(r)")
	s.root.printRecursive(" ")
}
