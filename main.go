package main

import (
	"fmt"

	"github.com/ororsatti/go-searchdex/search"
)

type testStruct struct {
	Id          string `index:"id"`
	Name        string `index:"text"`
	Description string `index:"text"`
}

type doc struct {
	Id      string `index:"id"`
	Content string `index:"text"`
}

func main() {
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
	index := search.New(anys)
	fmt.Println(index.Search("fox", 0))
}
