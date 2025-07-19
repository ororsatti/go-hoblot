package parsetags

import (
	"errors"
	"fmt"
	"reflect"
)

type indexTag string

const (
	tagName          = "index"
	idTag   indexTag = "id"
	textTag indexTag = "text"
)

type TagParser struct {
	original   any
	ref        reflect.Type
	idField    string
	textFields []string
}

func NewTagParser(o any) *TagParser {
	r := reflect.TypeOf(o)
	tp := &TagParser{
		original: o,
		ref:      r,
	}

	for i := range r.NumField() {
		fieldName := r.Field(i).Name
		indexType := indexTag(r.Field(i).Tag.Get(tagName))

		switch indexType {
		case idTag:
			{
				tp.idField = fieldName
			}

		case textTag:
			{
				tp.textFields = append(tp.textFields, fieldName)
			}
		}
	}
	return tp
}

func (t *TagParser) GetID() (string, error) {
	val, ok := t.ref.FieldByName(t.idField)
	if !ok {
		return "", errors.New("Missing ID field")
	}

	fmt.Println(val)
	return "", nil
}
