package model

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidGUID = errors.New("GUID is invalid")
)

type GUID struct {
	value string
}

func (g GUID) String() string {
	return g.value
}

func NewGUID(s string) (GUID, error) {
	r := regexp.MustCompile(`^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`)
	if !r.MatchString(s) {
		return GUID{value: ""}, ErrInvalidGUID
	}

	return GUID{value: s}, nil
}
