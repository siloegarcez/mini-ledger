package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const zeroBRLMoneyString = "0.00"

func TestNewBRLMoneyFromString(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantString string
		wantInt64  int64
		wantErr    bool
	}{
		{
			name:       "valid amount with 2 decimals",
			input:      "123.45",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "valid amount with 1 decimal",
			input:      "123.4",
			wantString: "123.40",
			wantInt64:  12340,
			wantErr:    false,
		},
		{
			name:       "valid amount with no decimals",
			input:      "123",
			wantString: "123.00",
			wantInt64:  12300,
			wantErr:    false,
		},
		{
			name:       "ivalid amount with dot but no decimals",
			input:      "123.",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "ivalid amount with dot but no integer part",
			input:      ".32",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "zero amount",
			input:      "0",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    false,
		},
		{
			name:       "zero with decimals",
			input:      zeroBRLMoneyString,
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    false,
		},
		{
			name:       "negative amount",
			input:      "-50.25",
			wantString: "-50.25",
			wantInt64:  -5025,
			wantErr:    false,
		},
		{
			name:       "small amount",
			input:      "0.01",
			wantString: "0.01",
			wantInt64:  1,
			wantErr:    false,
		},
		{
			name:       "large amount",
			input:      "999999.99",
			wantString: "999999.99",
			wantInt64:  99999999,
			wantErr:    false,
		},
		{
			name:       "amount with leading zeros",
			input:      "00123.45",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "too many decimals",
			input:      "123.456",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "invalid format (letters)",
			input:      "abc",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "double minus sign in integer part",
			input:      "--1",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "invalid decimal characters",
			input:      "1.a",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "negative decimal part",
			input:      "1.-1",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "invalid format (multiple dots)",
			input:      "123.45.67",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "empty string",
			input:      "",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "just a dot",
			input:      ".",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "whitespace handling",
			input:      "  123.45  ",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "invalid syntax with spaces",
			input:      "1 23.4 5",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "invalid syntax with spaces and multiple dots negative",
			input:      "-1. 23.4 5",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "two dots with negative sign",
			input:      "-123.-45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "dash inside integer part",
			input:      "-123-45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "explicit positive",
			input:      "+123.45",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "invalid characters in integer part",
			input:      "123-.-45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "negative zero",
			input:      "-0",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    false,
		},
		{
			name:       "negative zero with decimals",
			input:      "-" + zeroBRLMoneyString,
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    false,
		},
		{
			name:       "negative very small",
			input:      "-0.01",
			wantString: "-0.01",
			wantInt64:  -1,
			wantErr:    false,
		},
		{
			name:       "double positive sign",
			input:      "++123.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "plus and minus together",
			input:      "+-123.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "plus inside number",
			input:      "12+3.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "trailing characters",
			input:      "123.45abc",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "leading characters",
			input:      "abc123.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "currency prefix",
			input:      "R$123.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "comma thousand separator",
			input:      "1,234.56",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "dot thousand separator brazil style",
			input:      "1.234,56",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "max int64 cents boundary",
			input:      "92233720368547758.07",
			wantString: "92233720368547758.07",
			wantInt64:  9223372036854775807,
			wantErr:    false,
		},
		{
			name:       "overflow by one cent",
			input:      "92233720368547758.08",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "multiple leading zeros",
			input:      "000000.01",
			wantString: "0.01",
			wantInt64:  1,
			wantErr:    false,
		},
		{
			name:       "only zeros",
			input:      "0000",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    false,
		},
		{
			name:       "decimal with leading zero",
			input:      "0.10",
			wantString: "0.10",
			wantInt64:  10,
			wantErr:    false,
		},
		{
			name:       "decimal with trailing zero",
			input:      "10.10",
			wantString: "10.10",
			wantInt64:  1010,
			wantErr:    false,
		},
		{
			name:       "decimal exactly two zeros",
			input:      "10.00",
			wantString: "10.00",
			wantInt64:  1000,
			wantErr:    false,
		},
		{
			name:       "space between sign and number",
			input:      "- 123.45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "space inside decimal",
			input:      "123. 45",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "only spaces",
			input:      "   ",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "sign only",
			input:      "-",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "plus only",
			input:      "+",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
		{
			name:       "tab and newline handling",
			input:      "\t123.45\n",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "non-breaking space",
			input:      "123.45\u00A0",
			wantString: "123.45",
			wantInt64:  12345,
			wantErr:    false,
		},
		{
			name:       "zero-width space",
			input:      "123.45\u200B",
			wantString: zeroBRLMoneyString,
			wantInt64:  0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewBRLMoneyFromString(tt.input)

			if tt.wantErr {
				require.Error(t, err, "NewBRLMoney() expected error")
				return
			}

			require.NoError(t, err, "NewBRLMoney() unexpected error")
			assert.Equal(t, tt.wantInt64, got.Int64(), "NewBRLMoney() cents mismatch")
			assert.Equal(t, tt.wantString, got.String(), "NewBRLMoney() string mismatch")
		})
	}
}

func TestBRLMoney_Add(t *testing.T) { //nolint:dupl
	tests := []struct {
		name       string
		value      string
		otherValue string
		want       string
	}{
		{name: "adds two positive amounts", value: "10.25", otherValue: "2.75", want: "13.00"},
		{
			name:       "adds positive and negative amounts",
			value:      "10.00",
			otherValue: "-2.50",
			want:       "7.50",
		},
		{name: "adds two negative amounts", value: "-1.25", otherValue: "-0.75", want: "-2.00"},
		{name: "adds zero", value: "10.00", otherValue: "0", want: "10.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			m2, err := domain.NewBRLMoneyFromString(tt.otherValue)
			require.NoError(t, err)

			got := m.Add(m2)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestBRLMoney_Sub(t *testing.T) { //nolint:dupl
	tests := []struct {
		name       string
		value      string
		otherValue string
		want       string
	}{
		{
			name:       "subtracts smaller positive from larger positive",
			value:      "10.00",
			otherValue: "2.50",
			want:       "7.50",
		},
		{
			name:       "subtracts larger positive from smaller positive",
			value:      "2.50",
			otherValue: "10.00",
			want:       "-7.50",
		},
		{name: "subtracts negative", value: "10.00", otherValue: "-2.50", want: "12.50"},
		{name: "subtracts zero", value: "10.00", otherValue: "0", want: "10.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			m2, err := domain.NewBRLMoneyFromString(tt.otherValue)
			require.NoError(t, err)

			got := m.Sub(m2)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestBRLMoney_Neg(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "negates positive amount", value: "10.50", want: "-10.50"},
		{name: "negative becomes positive", value: "-10.50", want: "10.50"},
		{name: "keeps zero unchanged", value: "0", want: zeroBRLMoneyString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.Neg()
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestBRLMoney_IsNegative(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "positive amount is not negative", value: "1.00", want: false},
		{name: "negative amount is negative", value: "-1.00", want: true},
		{name: "zero is not negative", value: "0", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.IsNegative()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBRLMoney_IsPositive(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "positive amount is positive", value: "1.00", want: true},
		{name: "negative amount is not positive", value: "-1.00", want: false},
		{name: "zero is not positive", value: "0", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.IsPositive()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBRLMoney_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "positive amount is not zero", value: "1.00", want: false},
		{name: "negative amount is not zero", value: "-1.00", want: false},
		{name: "zero is zero", value: "0", want: true},
		{name: "negative zero parses as zero", value: "-0", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.IsZero()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBRLMoney_Abs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		value string
		want  string
	}{
		{name: "absolute value of positive amount is unchanged", value: "10.50", want: "10.50"},
		{name: "absolute value of negative amount is positive", value: "-10.50", want: "10.50"},
		{name: "absolute value of zero is zero", value: "0", want: zeroBRLMoneyString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.Abs()
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestBRLMoney_Int64(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int64
	}{
		{name: "zero returns zero cents", value: "0", want: 0},
		{name: "positive amount returns cents", value: "123.45", want: 12345},
		{name: "negative amount returns negative cents", value: "-123.45", want: -12345},
		{name: "single decimal is normalized to cents", value: "10.5", want: 1050},
		{
			name:  "max int64 cents boundary",
			value: "92233720368547758.07",
			want:  9223372036854775807,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.Int64()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBRLMoney_String(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "zero formats with 2 decimals", value: "0", want: zeroBRLMoneyString},
		{name: "positive amount stays normalized", value: "123.45", want: "123.45"},
		{name: "positive with one decimal is padded", value: "123.4", want: "123.40"},
		{name: "negative amount keeps sign", value: "-7.01", want: "-7.01"},
		{name: "negative zero formats as zero", value: "-0", want: zeroBRLMoneyString},
		{
			name:  "large value remains exact",
			value: "92233720368547758.07",
			want:  "92233720368547758.07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := domain.NewBRLMoneyFromString(tt.value)
			require.NoError(t, err)

			got := m.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
