package radix

import (
	"fmt"
	"iter"
	"maps"
	"sort"
	"strings"
)

type node[T any] struct {
	Children map[string]*node[T]
	Data     *T
}

type SearchResult struct {
	Distance int
	Data     any
}

type nodePath[T any] struct {
	node *node[T]
	key  string
}

func NewSearchResult(distance int, data any) SearchResult {
	return SearchResult{
		Distance: distance,
		Data:     data,
	}
}

func newNode[T any]() *node[T] {
	return &node[T]{
		Children: make(map[string]*node[T]),
	}
}

func (n *node[T]) addChild(key string, newNode *node[T]) {
	n.Children[key] = newNode
}

func (n *node[T]) removeChild(key string) {
	delete(n.Children, key)
}

func (n *node[T]) lookup(key string) *node[T] {
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

func cleanup[T any](path []nodePath[T]) {
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

func merge[T any](path []nodePath[T], keyToMerge string, nodeToMerge *node[T]) {
	if len(path) == 0 {
		return
	}

	last := path[len(path)-1]
	last.node.Children[last.key+keyToMerge] = nodeToMerge
	delete(last.node.Children, last.key)
}

func trackDown[T any](n *node[T], key string, path []nodePath[T]) (*node[T], []nodePath[T]) {
	if key == "" || n == nil {
		return n, path
	}

	for childKey, child := range n.Children {
		if strings.HasPrefix(key, childKey) {
			path = append(path, nodePath[T]{node: n, key: childKey})

			return trackDown(child, key[len(childKey):], path)
		}
	}

	path = append(path, nodePath[T]{node: n, key: key})
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

func (n *node[T]) createPath(key string) *node[T] {
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
				intermediate := newNode[T]()
				intermediate.addChild(childKey[splitIndex:], child)

				currentNode.addChild(childKey[:splitIndex], intermediate)
				currentNode.removeChild(childKey)

				currentNode = intermediate
			}

			continue LOOP_KEY
		}

		child := newNode[T]()
		currentNode.addChild(key[i:], child)
		return child
	}

	return currentNode
}

func (node *node[T]) fuzzyRecurse(
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

		// finding the max possible key len.
		// if the key length surpasses the amount of columns left, limit it to the columns left
		// and calculate the rest as a carry-over.
		maxKeyLen := min(len(key), n+maxDistance-m)
		carry := max(0, len(key)-maxKeyLen)

		for pos := range maxKeyLen {
			thisRowOffset := n * i
			prevRowOffset := thisRowOffset - n
			minDistance := matrix[thisRowOffset]

			for j, queryChar := range query {

				diff := convertBool(queryChar != rune(key[pos]))

				replaceCost := matrix[prevRowOffset+j] + diff
				deleteCost := matrix[prevRowOffset+j+1] + 1
				insertCost := matrix[thisRowOffset+j] + 1

				dist := min(replaceCost, deleteCost, insertCost)

				if dist < minDistance {
					minDistance = dist
				}

				matrix[thisRowOffset+j+1] = dist
			}

			if minDistance+carry > maxDistance {
				continue ITER_CHILDREN
			}

			i++
		}

		child.fuzzyRecurse(query, prefix+key, matrix, maxDistance, i, n, results)
	}
}

func (n *node[T]) printRecursive(prefix string) {
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

		fmt.Printf("%s%s %s\n", prefix, branch, key)

		newPrefix := prefix
		if isLastChild {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		child.printRecursive(newPrefix)
	}
}

func convertBool(b bool) int {
	if b {
		return 1
	}

	return 0
}
