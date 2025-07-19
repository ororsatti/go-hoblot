package main

import (
	parsetags "github.com/ororsatti/go-searchdex/parse_tags"
)

type testStruct struct {
	id          string `index:"id"`
	name        string `index:"text"`
	description string
}

func main() {
	test := testStruct{
		id:          "123",
		name:        "Jane",
		description: "A masterious woman",
	}

	tp := parsetags.NewTagParser(test)
	tp.GetID()
}
