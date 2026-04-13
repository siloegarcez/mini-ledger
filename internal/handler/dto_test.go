package handler_test

import (
	"mini-ledger/internal/handler"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumber_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    handler.Number
		wantErr bool
	}{
		{
			name:    "valid integer number",
			data:    []byte(`123`),
			want:    handler.Number("123"),
			wantErr: false,
		},
		{
			name:    "valid decimal number",
			data:    []byte(`123.45`),
			want:    handler.Number("123.45"),
			wantErr: false,
		},
		{
			name:    "json string number is accepted",
			data:    []byte(`"123.45"`),
			want:    handler.Number("123.45"),
			wantErr: false,
		},
		{
			name:    "invalid token",
			data:    []byte(`abc`),
			want:    handler.Number(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n handler.Number
			gotErr := n.UnmarshalJSON(tt.data)
			if gotErr != nil {
				assert.True(t, tt.wantErr)
				return
			}
			require.False(t, tt.wantErr)
			assert.Equal(t, tt.want, n)
		})
	}
}

func TestNumber_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       handler.Number
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal decimal number",
			n:       handler.Number("123.45"),
			want:    []byte(`123.45`),
			wantErr: false,
		},
		{
			name:    "marshal integer number",
			n:       handler.Number("123"),
			want:    []byte(`123`),
			wantErr: false,
		},
		{
			name:    "invalid number format returns error",
			n:       handler.Number("12-3"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.n.MarshalJSON()
			if gotErr != nil {
				assert.True(t, tt.wantErr)
				return
			}
			require.False(t, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
