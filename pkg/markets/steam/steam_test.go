package steam

import (
	"reflect"
	"testing"

	"github.com/arseniizyk/investor1337/pkg/markets"
)

func Test_format(t *testing.T) {
	tests := []struct {
		name     string
		input    *Response
		expected []markets.Pair
		wantErr  bool
	}{
		{
			name: "valid with comma price",
			input: &Response{
				SellOrderTable: []struct {
					Price        string `json:"price"`
					PriceWithFee string `json:"price_with_fee"`
					Quantity     string `json:"quantity"`
				}{
					{Price: "1,234.56", Quantity: "10"},
				},
			},
			expected: []markets.Pair{{Price: 1234.56, Quantity: 10}},
		},
		{
			name: "invalid price",
			input: &Response{
				SellOrderTable: []struct {
					Price        string `json:"price"`
					PriceWithFee string `json:"price_with_fee"`
					Quantity     string `json:"quantity"`
				}{
					{Price: "NaN", Quantity: "5"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid quantity",
			input: &Response{
				SellOrderTable: []struct {
					Price        string `json:"price"`
					PriceWithFee string `json:"price_with_fee"`
					Quantity     string `json:"quantity"`
				}{
					{Price: "12.3", Quantity: "XX"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty list",
			input: &Response{SellOrderTable: []struct {
				Price        string `json:"price"`
				PriceWithFee string `json:"price_with_fee"`
				Quantity     string `json:"quantity"`
			}{}},
			expected: []markets.Pair{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := format(tt.input)
			if err == nil && tt.wantErr {
				t.Fatalf("format() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("format() got = %+v, expected %+v", got, tt.expected)
			}
		})
	}
}
