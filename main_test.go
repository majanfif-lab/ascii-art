package main

import (
	"os"
	"strings"
	"testing"
)

// Helper function to load map for tests
func loadAsciiMap(t *testing.T) map[rune][]string {
	banner, err := os.ReadFile("standard.txt")
	if err != nil {
		t.Fatalf("Could not read standard.txt for testing: %v", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(banner), "\r", ""), "\n")
	asciiMap := make(map[rune][]string)
	for i := 32; i <= 126; i++ {
		start := (i-32)*9 + 1
		asciiMap[rune(i)] = lines[start : start+8]
	}
	return asciiMap
}

func TestGenerateAscii(t *testing.T) {
	asciiMap := loadAsciiMap(t)

	tests := []struct {
		name          string
		input         string
		expectedLines int // Number of expected \n returns.
		expectError   bool
	}{
		{
			name:          "Empty string",
			input:         "",
			expectedLines: 0,
		},
		{
			name:          "Single literal newline",
			input:         "\\n",
			expectedLines: 1,
		},
		{
			name:          "Multiple literal newlines",
			input:         "\\n\\n\\n",
			expectedLines: 3,
		},
		{
			name:          "Standard Word",
			input:         "hello",
			expectedLines: 8, // 8 lines of art
		},
		{
			name:          "Word with a trailing newline",
			input:         "Hello\\n",
			expectedLines: 9, // 8 lines for 'Hello' + 1 for \n
		},
		{
			name:          "Words separated by newlines",
			input:         "Hello\\n\\nThere",
			expectedLines: 17, // 8 for 'Hello' + 1 empty line + 8 for 'There'
		},
		{
			name:          "Invalid character",
			input:         "Hélló",
			expectedLines: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateAscii(tt.input, asciiMap)

			if tt.expectError && err == nil {
				t.Errorf("Expected an error for input %q, but got none", tt.input)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Did not expect error for input %q, but got: %v", tt.input, err)
			}

			// Count line returns to assert correct output length
			if !tt.expectError {
				actualLines := strings.Count(result, "\n")
				if actualLines != tt.expectedLines {
					t.Errorf("Input %q: expected %d newlines in output, got %d", tt.input, tt.expectedLines, actualLines)
				}
			}
		})
	}
}
