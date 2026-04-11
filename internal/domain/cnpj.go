package domain

import (
	"errors"
	"fmt"
	"strings"
)

type CNPJ struct {
	value string
}

const (
	cnpjLength = 14
)

var (
	firstDigitWeights  = [12]int32{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	secondDigitWeights = [13]int32{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

var (
	ErrInvalidCNPJ       = errors.New("invalid")
	ErrInvalidCNPJLength = errors.New("invalid cnpj length, expected 14 digits")
)

// Expected input format is "12345678000195" or "12.345.678/0001-95".
func ParseCNPJ(s string) (*CNPJ, error) {
	dvValues := []int32{}
	s = strings.ToUpper(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "-", ""),
				".",
				"",
			),
			"/",
			"",
		),
	)
	var b strings.Builder
	b.Grow(cnpjLength)
	for _, c := range s {
		if c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' {
			if b.Len()+1 > cnpjLength {
				return nil, ErrInvalidCNPJLength
			}
			b.WriteRune(c)
			dvValues = append(dvValues, c-'0')
		}
	}

	_cnpj := b.String()

	if len(_cnpj) != cnpjLength || len(dvValues) != cnpjLength {
		return nil, ErrInvalidCNPJ
	}

	firstDigitSum := int32(0)
	firstDigit := int32(0)
	for i := range firstDigitWeights {
		firstDigitSum += firstDigitWeights[i] * dvValues[i]
	}

	r := firstDigitSum % 11

	if r >= 2 {
		firstDigit = 11 - r
	}

	if firstDigit != dvValues[12] {
		return nil, fmt.Errorf("cnpj: %w first digit", ErrInvalidCNPJ)
	}

	secondDigitSum := int32(0)
	secondDigit := int32(0)
	for i := range secondDigitWeights {
		secondDigitSum += secondDigitWeights[i] * dvValues[i]
	}

	r = secondDigitSum % 11

	if r >= 2 {
		secondDigit = 11 - r
	}

	if secondDigit != dvValues[13] {
		return nil, fmt.Errorf("cnpj: %w second digit", ErrInvalidCNPJ)
	}

	return &CNPJ{value: _cnpj}, nil
}

// Format: "12345678000195".
func (c CNPJ) String() string {
	return c.value
}

// Format: "12.345.678/0001-95".
func (c CNPJ) Formatted() string {
	return c.value[0:2] + "." + c.value[2:5] + "." + c.value[5:8] + "/" + c.value[8:12] + "-" + c.value[12:14]
}
