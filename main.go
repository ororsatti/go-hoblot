package main

import (
	"fmt"
	"iter"
	"maps"
	"sort"
	"strings"

	"github.com/jdkato/prose/v2"
)

func convertBool(b bool) int {
	if b {
		return 1
	}

	return 0
}

type Node struct {
	Children map[string]*Node
	Data     any
}

type SearchableMap struct {
	root *Node
}

func createNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

func (n *Node) Size() int {
	return len(n.Children)
}

func (n *Node) addChild(key string, newNode *Node) {
	n.Children[key] = newNode
}

func (n *Node) removeChild(key string) {
	delete(n.Children, key)
}

func NewSearchableMap() *SearchableMap {
	return &SearchableMap{
		root: createNode(),
	}
}

func (s *SearchableMap) Set(key string) *Node {
	node := s.root.createPath(key)
	node.Data = true
	return node
}

func (s *SearchableMap) Get(key string) *Node {
	node := s.root.lookup(key)
	if node.Data != nil {
		return node
	}

	return nil
}

type SearchResult struct {
	Distance int
	Data     any
}

func NewSearchResult(distance int, data any) SearchResult {
	return SearchResult{
		Distance: distance,
		Data:     data,
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

func (node *Node) fuzzyRecurse(
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

type nodePath struct {
	node *Node
	key  string
}

func (n *Node) lookup(key string) *Node {
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

func merge(path []nodePath, keyToMerge string, nodeToMerge *Node) {
	if len(path) == 0 {
		return
	}

	last := path[len(path)-1]
	last.node.Children[last.key+keyToMerge] = nodeToMerge
	delete(last.node.Children, last.key)
}

func trackDown(n *Node, key string, path []nodePath) (*Node, []nodePath) {
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

func (n *Node) createPath(key string) *Node {
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
				intermediate := createNode()
				intermediate.addChild(childKey[splitIndex:], child)

				currentNode.addChild(childKey[:splitIndex], intermediate)
				currentNode.removeChild(childKey)

				currentNode = intermediate
			}

			continue LOOP_KEY
		}

		child := createNode()
		currentNode.addChild(key[i:], child)
		return child
	}

	return currentNode
}

func (n *Node) printRecursive(prefix string) {
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

func (s *SearchableMap) Print() {
	fmt.Println("(r)")
	s.root.printRecursive(" ")
}

func main() {
	text := `
	The most saddening thought that arises after the perusal of this Volume,
	is, that no change has yet been made in the infamous Lunacy Laws, for
	which, in the main, we have to thank our Whig Rulers. Never was a more
	criminal or despotic Law passed than that which now enables a Husband to
	lock up his Wife in a Madhouse on the certificate of two medical men,
	who often in haste, frequently for a bribe, certify to madness where
	none exists. We believe that under these Statutes thousands of persons,
	perfectly sane, are now imprisoned in private asylums throughout the
	Kingdom; while strangers are in possession of their property; and the
	miserable prisoner is finally brought to a state of actual lunacy or
	imbecility—however rational he may have been when first immured. The
	Keepers of these Madhouse Dens, from long study in their diabolical art,
	can reduce, by certain drugs, the clearest brain to a state of stupor;
	and the Lunacy Commissioners take all for granted that they hear over the
	luxurious lunch with which the Mad Doctor regales them.
	`

	doc, _ := prose.NewDocument(text)

	smap := NewSearchableMap()
	// Tokenize the text
	for _, tok := range doc.Tokens() {
		smap.Set(strings.ToLower(tok.Text))
	}

	res := smap.FuzzyGet("king", 1)
	fmt.Println(res)
	// for key, val := range res {
	// 	fmt.Println(key, val)
	// }
}
