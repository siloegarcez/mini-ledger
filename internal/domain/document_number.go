package domain

import (
	"errors"
	"fmt"
)

type DocumentNumber struct {
	value     string
	formatted string
	docType   DocumentType
}

var (
	ErrInvalidDocumentNumberLen = errors.New(
		"invalid document number length",
	)
	ErrInvalidDocumentNumber = errors.New(
		"invalid document number",
	)
)

const (
	MaxDocumentNumberLength = 18
)

type DocumentType int

const (
	DocumentTypeCPF DocumentType = iota + 1
	DocumentTypeCNPJ
)

func (dt DocumentType) String() string {
	switch dt {
	case DocumentTypeCPF:
		return "CPF"
	case DocumentTypeCNPJ:
		return "CNPJ"
	}
	return ""
}

func NewDocumentNumber(value string) (DocumentNumber, error) {
	if len(value) == 0 || len(value) > MaxDocumentNumberLength {
		return DocumentNumber{}, ErrInvalidDocumentNumberLen
	}

	cpf, err := ParseCPF(value)

	if err == nil {
		return DocumentNumber{
			value:     cpf.String(),
			formatted: cpf.Formatted(),
			docType:   DocumentTypeCPF,
		}, nil
	}

	cnpj, err := ParseCNPJ(value)

	if err == nil {
		return DocumentNumber{
			value:     cnpj.String(),
			formatted: cnpj.Formatted(),
			docType:   DocumentTypeCNPJ,
		}, nil
	}

	return DocumentNumber{}, ErrInvalidDocumentNumber
}

func (d DocumentNumber) String() string {
	return d.value
}

func (d DocumentNumber) Formatted() string {
	return d.formatted
}

func (d DocumentNumber) Type() DocumentType {
	return d.docType
}

func (d DocumentNumber) IsEmpty() bool {
	return len(d.value) == 0
}

func (d DocumentNumber) Validate() []error {
	errs := []error{}
	if d.IsEmpty() {
		errs = append(
			errs,
			fmt.Errorf("%w: document number cannot be empty", ErrInvalidDocumentNumber),
		)
	}
	if len(d.docType.String()) == 0 {
		errs = append(
			errs,
			fmt.Errorf("%w: document type cannot be empty", ErrInvalidDocumentNumber),
		)
	}
	if len(d.formatted) == 0 {
		errs = append(
			errs,
			fmt.Errorf("%w: formatted document number cannot be empty", ErrInvalidDocumentNumber),
		)
	}

	return errs
}
