package cli

import (
	"math/big"
	"testing"
)

func TestFormatEth(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{"Nil", nil, "0"},
		{"Zero", big.NewInt(0), "0"},
		{"1 Wei", big.NewInt(1), "0"},                 // 1e-18, formatted to 9 decimals is 0
		{"1 Gwei", big.NewInt(1e9), "0.000000001"},    // 1e-9
		{"1 ETH", big.NewInt(1e18), "1"},
		{"1.5 ETH", big.NewInt(1500000000000000000), "1.5"},
		{"Small", big.NewInt(100), "0"},               // 100 wei is 1e-16, formatted to 9 decimals is 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatEth(tt.input)
			if got != tt.expected {
				t.Errorf("formatEth(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatGwei(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{"Nil", nil, "0"},
		{"Zero", big.NewInt(0), "0"},
		{"1 Wei", big.NewInt(1), "0.000000001"},
		{"1 Gwei", big.NewInt(1e9), "1"},
		{"1.5 Gwei", big.NewInt(1500000000), "1.5"},
		{"0.1 Gwei", big.NewInt(100000000), "0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatGwei(tt.input)
			if got != tt.expected {
				t.Errorf("formatGwei(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Small", 123, "123"},
		{"Thousand", 1000, "1,000"},
		{"Million", 1234567, "1,234,567"},
		{"Billion", 1000000000, "1,000,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatUint64(tt.input)
			if got != tt.expected {
				t.Errorf("formatUint64(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
