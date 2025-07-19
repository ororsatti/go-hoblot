package parsetags

import (
	"errors"
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
	idField    string
	textFields []string
	refVal     reflect.Value
	invalid    bool
}

func NewTagParser(o any) *TagParser {
	r := reflect.TypeOf(o)
	tp := &TagParser{
		original: o,
		refVal:   reflect.ValueOf(o),
	}

	if r.Kind() != reflect.Struct {
		tp.invalid = true
		return tp
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
	if t.invalid {
		return "", errors.New("acting on invalid tag parser")
	}

	idVal := t.refVal.FieldByName(t.idField)
	if idVal.Kind() != reflect.String {
		return "", errors.New("id must be string")
	}

	return idVal.String(), nil
}

func (t *TagParser) GetText() ([]string, error) {
	if t.invalid {
		return nil, errors.New("acting on invalid tag parser")
	}

	var vals []string

	for _, field := range t.textFields {
		txtVal := t.refVal.FieldByName(field)
		if txtVal.Kind() != reflect.String {
			return nil, errors.New("all text fields must be string")
		}

		vals = append(vals, txtVal.String())

	}

	return vals, nil
}
