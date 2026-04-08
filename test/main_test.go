package test

import (
	"sewz/core"
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"long string", "hello world", 5, "he..."},
		{"zero max", "hello", 0, "hello"},
		{"negative max", "hello", -1, "hello"},
		{"single char", "h", 5, "h"},
		{"truncate long", "verylongstring", 10, "verylon..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.Truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"single port", "8080->8080/tcp", "8080->8080/tcp"},
		{"multiple ports", "8080->8080/tcp, 9090->9090/tcp", "8080->8080/tcp\n9090->9090/tcp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.FormatPorts(tt.input)
			if result != tt.expected {
				t.Errorf("FormatPorts(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
