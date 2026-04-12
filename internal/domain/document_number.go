package domain

import (
	"errors"
	"strings"
	"unicode"
)

type DocumentNumber struct {
	value string
}

var (
	ErrDocumentNumberInvalidLen = errors.New(
		"invalid document number length",
	)
	ErrDocumentNumberEmpty = errors.New(
		"document number cannot be empty",
	)
)

const (
	MaxDocumentNumberLength = 15
)

type DocumentType int

func NewDocumentNumber(value string) (DocumentNumber, error) {
	value = strings.ToLower(removeAllWhitespace(value))
	if len(value) == 0 {
		return DocumentNumber{}, ErrDocumentNumberEmpty
	}

	if len(value) > MaxDocumentNumberLength {
		return DocumentNumber{}, ErrDocumentNumberInvalidLen
	}

	return DocumentNumber{value: value}, nil
}

func (d DocumentNumber) String() string {
	return d.value
}

func removeAllWhitespace(s string) string {
	var b strings.Builder
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
