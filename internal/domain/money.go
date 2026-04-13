package domain

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BRLMoney struct {
	cents int64
	scale int16
}

var (
	ErrBRLMoneyAmountEmpty           = errors.New("empty money amount")
	ErrBRLMoneyInvalidDecimalPart    = errors.New("invalid decimal part")
	ErrBRLMoneyInvalidDecimalPartLen = errors.New("invalid decimal part length")
	ErrBRLMoneyInvalidIntegerPart    = errors.New("invalid integer part")
	ErrBRLMoneyInvalidAmount         = errors.New("invalid amount")
	ErrBRLMoneyAmountOverflow        = errors.New("money amount overflow")
)

func NewBRLMoneyFromString(value string) (BRLMoney, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return BRLMoney{}, ErrBRLMoneyAmountEmpty
	}

	negative := false
	if strings.HasPrefix(value, "-") {
		negative = true
		value = value[1:]
	}

	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return BRLMoney{}, ErrBRLMoneyInvalidAmount
	}

	if strings.Contains(value, ".") && (len(parts[0]) == 0 || len(parts[1]) == 0) {
		return BRLMoney{}, ErrBRLMoneyInvalidAmount
	}

	integerPartStr := parts[0]
	decimalPartStr := ""

	if len(parts) == 2 {
		decimalPartStr = parts[1]
	}

	if len(decimalPartStr) > 2 {
		return BRLMoney{}, ErrBRLMoneyInvalidDecimalPartLen
	}

	for len(decimalPartStr) < 2 {
		decimalPartStr += "0"
	}

	intPart, err := strconv.ParseInt(integerPartStr, 10, 64)
	if err != nil {
		return BRLMoney{}, ErrBRLMoneyInvalidIntegerPart
	}
	if intPart < 0 {
		return BRLMoney{}, ErrBRLMoneyInvalidIntegerPart
	}

	decPart, err := strconv.ParseInt(decimalPartStr, 10, 64)
	if err != nil {
		return BRLMoney{}, ErrBRLMoneyInvalidDecimalPart
	}
	if decPart < 0 {
		return BRLMoney{}, ErrBRLMoneyInvalidDecimalPart
	}

	if intPart > (math.MaxInt64-decPart)/100 {
		return BRLMoney{}, ErrBRLMoneyAmountOverflow
	}

	amount := intPart*100 + decPart

	if negative {
		amount = -amount
	}

	return BRLMoney{
		cents: amount,
		scale: 2,
	}, nil
}

func (m BRLMoney) Int64() int64 {
	return m.cents
}

func (m BRLMoney) String() string {
	sign := ""
	amount := uint64(m.cents) //nolint:gosec

	if m.cents < 0 {
		sign = "-"
		amount = uint64(^m.cents) + 1 // nolint:gosec
	}

	intPart := amount / 100
	decPart := amount % 100

	return fmt.Sprintf("%s%d.%02d", sign, intPart, decPart)
}

func (m BRLMoney) Add(m2 BRLMoney) BRLMoney {
	return BRLMoney{
		cents: m.cents + m2.cents,
		scale: m.scale,
	}
}

func (m BRLMoney) Sub(m2 BRLMoney) BRLMoney {
	return BRLMoney{
		cents: m.cents - m2.cents,
		scale: m.scale,
	}
}

func (m BRLMoney) Neg() BRLMoney {
	if m.IsZero() {
		return m
	}
	return BRLMoney{
		cents: -m.cents,
		scale: m.scale,
	}
}

func (m BRLMoney) Abs() BRLMoney {
	if m.IsZero() || m.IsPositive() {
		return m
	}
	return BRLMoney{
		cents: -m.cents,
		scale: m.scale,
	}
}

func (m BRLMoney) IsNegative() bool {
	return m.cents < 0
}

func (m BRLMoney) IsPositive() bool {
	return m.cents > 0
}

func (m BRLMoney) IsZero() bool {
	return m.cents == 0
}
