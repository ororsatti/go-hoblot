package radix

import (
	"fmt"
	"io"
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
	results := make(map[string]SearchResult)

	w := len(query) + 1
	h := w + maxDistance

	mat := make([][]int, h)
	for i := range mat {
		mat[i] = make([]int, w)
	}

	for i := range w {
		mat[0][i] = i
	}

	for i := range h {
		mat[i][0] = i
	}

	s.root.fuzzyRecurse2(query, "", mat, maxDistance, 0, results)

	return results
}

func (s *SearchableMap) Print(out io.Writer) {
	fmt.Fprintln(out, "(r)")
	s.root.printRecursive(out, " ")
}
