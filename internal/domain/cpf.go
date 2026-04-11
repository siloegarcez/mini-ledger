package domain

import (
	"errors"
	"strings"
)

type CPF struct {
	value string
}

var (
	blacklist = map[string]struct{}{
		"00000000000": {},
		"11111111111": {},
		"22222222222": {},
		"33333333333": {},
		"44444444444": {},
		"55555555555": {},
		"66666666666": {},
		"77777777777": {},
		"88888888888": {},
		"99999999999": {},
		"12345678909": {},
	}
)

var (
	ErrInvalidCPF       = errors.New("invalid cpf")
	ErrInvalidCPFLength = errors.New("invalid cpf length, expected 11 digits")
)

const (
	cpfLength = 11
)

// Format: "12345678909".
func (c CPF) String() string {
	return c.value
}

// Format: "123.456.789-09".
func (c CPF) Formatted() string {
	return c.value[0:3] + "." + c.value[3:6] + "." + c.value[6:9] + "-" + c.value[9:11]
}

// Expected input format is "12345678909" or "123.456.789-09".
func ParseCPF(s string) (*CPF, error) {
	s = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "-", ""), ".", "")
	var b strings.Builder
	b.Grow(cpfLength)
	for _, c := range s {
		if c >= '0' && c <= '9' {
			if b.Len()+1 > cpfLength {
				return nil, ErrInvalidCPFLength
			}
			b.WriteRune(c)
		}
	}

	_cpf := b.String()

	if len(_cpf) != cpfLength {
		return nil, ErrInvalidCPFLength
	}

	if _, ok := blacklist[_cpf]; ok {
		return nil, ErrInvalidCPF
	}

	firstCalcDigit := 0
	k := 11
	for _, char := range _cpf[:9] {
		k--
		digit := int(char - '0')
		firstCalcDigit += digit * k
	}
	firstDigit := int(_cpf[9] - '0')
	firstCalcDigit = firstCalcDigit * 10 % 11

	if firstCalcDigit == 10 {
		firstCalcDigit = 0
	}

	if firstDigit != firstCalcDigit {
		return nil, ErrInvalidCPF
	}

	k = 12
	secondCalcDigit := 0
	for _, char := range _cpf[:10] {
		k--
		digit := int(char - '0')
		secondCalcDigit += digit * k
	}
	secondCalcDigit = secondCalcDigit * 10 % 11

	if secondCalcDigit == 10 {
		secondCalcDigit = 0
	}

	secondDigit := int(_cpf[10] - '0')

	if secondDigit != secondCalcDigit {
		return nil, ErrInvalidCPF
	}

	return &CPF{value: _cpf}, nil
}
