package parsetags

import "testing"

type testStruct struct {
	id          string
	name        string
	description string
}

func TestGetID(t *testing.T) {
	test := testStruct{
		id:          "123",
		name:        "Jane",
		description: "A masterious woman",
	}

	tp := NewTagParser(test)
	tp.GetID()
}
