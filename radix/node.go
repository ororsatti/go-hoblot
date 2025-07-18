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
	for i := 0; ; i++ {
		if len(src) <= i {
			return len(src)
		}

		if len(dest) <= i {
			return len(dest)
		}

		if dest[i] != src[i] {
			return i
		}
	}
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

func (node *node) fuzzyRecurse(
	query, prefix string,
	matrix []int,
	maxDistance, m, n int,
	results map[string]SearchResult,
) {
	if node.Data != nil {
		if matrix[m*n-1] <= maxDistance {
			results[prefix] = NewSearchResult(matrix[m*n-1], node.Data)
		}
	}

ITER_CHILDREN:
	for key, child := range node.Children {

		i := m

		for pos := range len(key) {

			thisRowOffset := n * i
			prevRowOffset := thisRowOffset - n

			minDistance := matrix[thisRowOffset]

			for j := 0; j < n-1; j++ {

				diff := convertBool(query[j] != key[pos])

				replaceCost := matrix[prevRowOffset+j] + diff
				deleteCost := matrix[prevRowOffset+j+1] + 1
				insertCost := matrix[thisRowOffset+j] + 1

				dist := min(replaceCost, deleteCost, insertCost)

				if dist < minDistance {
					minDistance = dist
				}

				matrix[thisRowOffset+j+1] = dist
			}

			if minDistance > maxDistance {
				continue ITER_CHILDREN
			}

			i++
		}

		child.fuzzyRecurse(query, prefix+key, matrix, maxDistance, i, n, results)
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

		fmt.Fprintf(out, "%s%s %s\n", prefix, branch, key)

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
