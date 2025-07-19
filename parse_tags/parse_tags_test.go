package parsetags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID          string `index:"id"`
	Name        string `index:"text"`
	Description string
}

func TestGetID(t *testing.T) {
	testCases := []struct {
		name string
		u    User
		id   string
	}{
		{
			name: "empty",
			u:    User{},
			id:   "",
		},
		{
			name: "with id",
			u: User{
				ID:          "123",
				Name:        "Jane",
				Description: "A masterious woman",
			},
			id: "123",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tp := NewTagParser(testCase.u)

			id, err := tp.GetID()
			assert.NoError(tt, err)
			assert.Equal(tt, testCase.id, id)
		})
	}
}

func TestGetID_Err(t *testing.T) {
	tp := NewTagParser(true)

	_, err := tp.GetID()
	assert.Error(t, err)
}

func TestGetText(t *testing.T) {
	testCases := []struct {
		u    User
		text []string
		name string
	}{
		{
			name: "empty",
			u:    User{},
			text: []string{""},
		},
		{
			name: "with text",
			u: User{
				ID:          "123",
				Name:        "Jane",
				Description: "A masterious woman",
			},
			text: []string{"Jane"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tp := NewTagParser(testCase.u)

			text, err := tp.GetText()
			assert.NoError(tt, err)
			assert.Equal(tt, testCase.text, text)
		})
	}
}

func TestGetText_Err(t *testing.T) {
	tp := NewTagParser(true)

	_, err := tp.GetText()
	assert.Error(t, err)
}