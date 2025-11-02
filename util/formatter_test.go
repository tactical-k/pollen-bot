package util

import "testing"

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard date format",
			input:    "2024-03-15",
			expected: "2024年3月15日",
		},
		{
			name:     "January date",
			input:    "2024-01-01",
			expected: "2024年1月1日",
		},
		{
			name:     "December date",
			input:    "2024-12-31",
			expected: "2024年12月31日",
		},
		{
			name:     "Invalid format - too few parts",
			input:    "2024-03",
			expected: "2024-03",
		},
		{
			name:     "Invalid format - too many parts",
			input:    "2024-03-15-extra",
			expected: "2024-03-15-extra",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No dashes",
			input:    "20240315",
			expected: "20240315",
		},
		{
			name:     "Different separator",
			input:    "2024/03/15",
			expected: "2024/03/15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.input)
			if result != tt.expected {
				t.Errorf("FormatDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
