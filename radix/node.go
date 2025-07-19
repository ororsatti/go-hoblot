package radix

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"sort"
	"strings"
)

type node struct {
	Children map[string]*node
	Data     any
}

type SearchResult struct {
	Distance int
	Data     any
}

type nodePath struct {
	node *node
	key  string
}

func NewSearchResult(distance int, data any) SearchResult {
	return SearchResult{
		Distance: distance,
		Data:     data,
	}
}

func newNode() *node {
	return &node{
		Children: make(map[string]*node),
	}
}

func (n *node) addChild(key string, newNode *node) {
	n.Children[key] = newNode
}

func (n *node) removeChild(key string) {
	delete(n.Children, key)
}

func (n *node) lookup(key string) *node {
	if key == "" {
		return n
	}

	for childKey, child := range n.Children {
		if strings.HasPrefix(key, childKey) {
			return child.lookup(key[len(childKey):])
		}
	}

	return n
}

func cleanup(path []nodePath) {
	if len(path) == 0 {
		return
	}

	p := path[len(path)-1]

	delete(p.node.Children, p.key)

	if p.node.Data != nil {
		return
	}

	next, stop := iter.Pull2(maps.All(p.node.Children))
	defer stop()

	firstKey, firstChild, ok := next()
	newPath := path[:len(path)-1]

	if !ok {
		cleanup(newPath)
	} else {
		merge(newPath, firstKey, firstChild)
	}
}

func merge(path []nodePath, keyToMerge string, nodeToMerge *node) {
	if len(path) == 0 {
		return
	}

	last := path[len(path)-1]
	last.node.Children[last.key+keyToMerge] = nodeToMerge
	delete(last.node.Children, last.key)
}

func trackDown(n *node, key string, path []nodePath) (*node, []nodePath) {
	if key == "" || n == nil {
		return n, path
	}

	for childKey, child := range n.Children {
		if strings.HasPrefix(key, childKey) {
			path = append(path, nodePath{node: n, key: childKey})

			return trackDown(child, key[len(childKey):], path)
		}
	}

	path = append(path, nodePath{node: n, key: key})
	return trackDown(nil, "", path)
}

func findSplitIndex(src string, dest string) int {
	srcRunes := []rune(src)
	destRunes := []rune(dest)

	lenSrc := len(srcRunes)
	lenDest := len(destRunes)

	minLen := lenSrc
	if lenDest < minLen {
		minLen = lenDest
	}

	for i := 0; i < minLen; i++ {
		if srcRunes[i] != destRunes[i] {
			return i
		}
	}

	return minLen
}

func (n *node) createPath(key string) *node {
	currentNode := n

LOOP_KEY:
	for i := 0; i < len(key); {
		partialKey := key[i:]

		for childKey, child := range currentNode.Children {
			splitIndex := findSplitIndex(partialKey, childKey)
			if splitIndex == 0 {
				continue
			}

			i += splitIndex

			if splitIndex == len(childKey) {
				currentNode = child
			} else {
				intermediate := newNode()
				intermediate.addChild(childKey[splitIndex:], child)

				currentNode.addChild(childKey[:splitIndex], intermediate)
				currentNode.removeChild(childKey)

				currentNode = intermediate
			}

			continue LOOP_KEY
		}

		child := newNode()
		currentNode.addChild(key[i:], child)
		return child
	}

	return currentNode
}

func (n *node) fuzzyRecurse(query string, prefix string,
	matrix [][]int,
	maxDistance int,
	row int,
	results map[string]SearchResult,
) {
	if n.Data != nil {
		dist := matrix[row][len(query)]

		if dist <= maxDistance {
			results[prefix] = NewSearchResult(dist, n.Data)
		}
	}

	maxRows := len(matrix)
	if row >= maxRows {
		return
	}

	queryRunes := []rune(query)

KEY_ITER:
	for key, child := range n.Children {
		keyRunes := []rune(key)
		m := row

		// fmt.Printf(">>>> %s <<<<\n", key)
		for _, keyRune := range keyRunes {
			m++

			if m >= maxRows {
				continue KEY_ITER
			}

			prevRow := matrix[m-1]
			currRow := matrix[m]

			for j, qRune := range queryRunes {
				cost := 0
				if keyRune != qRune {
					cost = 1
				}

				currRow[j+1] = min(
					prevRow[j]+cost, // substitution/match
					currRow[j]+1,    // insertion
					prevRow[j+1]+1,  // deletion
				)
			}

		}

		if m < len(query)+maxDistance+1 {
			child.fuzzyRecurse(query, prefix+key, matrix, maxDistance, m, results)
		}

		// fmt.Printf(">>>> end %s <<<<\n", key)
	}
}

func (n *node) printRecursive(out io.Writer, prefix string) {
	keys := make([]string, 0, len(n.Children))
	for key := range n.Children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for i, key := range keys {
		child := n.Children[key]
		isLastChild := i == len(keys)-1

		branch := "├──"
		if isLastChild {
			branch = "└──"
		}

		mark := ""
		if child.Data != nil {
			mark = "*"
		}

		fmt.Fprintf(out, "%s%s %s %s\n", prefix, branch, key, mark)

		newPrefix := prefix
		if isLastChild {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		child.printRecursive(out, newPrefix)
	}
}

func convertBool(b bool) int {
	if b {
		return 1
	}

	return 0
}
