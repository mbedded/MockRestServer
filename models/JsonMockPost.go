package models

import "strings"

type JsonMockPost struct {
	Key     string
	Content string
}

func (mock *JsonMockPost) TrimFields() {
	mock.Key = strings.TrimSpace(mock.Key)
	mock.Content = strings.TrimSpace(mock.Content)
}
