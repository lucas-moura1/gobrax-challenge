package entity

import "strings"

type ErrorInvalidField struct {
	Message []string
}

func (e ErrorInvalidField) Error() string {
	return strings.Join(e.Message, ",")
}
