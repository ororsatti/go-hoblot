package radix

import (
	"fmt"
	"iter"
	"maps"
)

type SearchableMap[T any] struct {
	root *node[T]
}

func NewSearchableMap[T any]() *SearchableMap[T] {
	return &SearchableMap[T]{
		root: newNode[T](),
	}
}

func (s *SearchableMap[T]) Set(key string, dataPtr *T) *node[T] {
	node := s.root.createPath(key)
	node.Data = dataPtr
	return node
}

func (s *SearchableMap[T]) Get(key string) *node[T] {
	node := s.root.lookup(key)
	if node.Data != nil {
		return node
	}

	return nil
}

func (s *SearchableMap[T]) Delete(key string) {
	node, path := trackDown(s.root, key, []nodePath[T]{})
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

func (s *SearchableMap[T]) FuzzyGet(query string, maxDistance int) map[string]SearchResult {
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

func (s *SearchableMap[T]) Print() {
	fmt.Println("(r)")
	s.root.printRecursive(" ")
}
