package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// 1. Check argument length

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER (optional)]")
		return
	}

	input := os.Args[1]
	if input == "" {
		return
	}

	// Default banner file is standard.txt
	bannerFile := "standard.txt"

	// If you type a 2nd argument (like "shadow"), use that file instead!
	if len(os.Args) == 3 {
		bannerFile = os.Args[2] + ".txt"
	}

	// Read the chosen banner file
	banner, err := os.ReadFile(bannerFile)
	if err != nil {
		fmt.Printf("Error: could not read banner file %s\n", bannerFile)
		return
	}
	// ... (rest of your code stays the same) ...

	// Normalize newlines (in case of Windows CRLF formats) and split into lines
	lines := strings.Split(strings.ReplaceAll(string(banner), "\r", ""), "\n")

	if len(lines) < 855 {
		fmt.Println("Error: standard.txt is corrupted or invalid.")
		return
	}

	// 3. Map the ASCII characters (Printable ASCII starts at 32 ' ' up to 126 '~')
	asciiMap := make(map[rune][]string)
	for i := 32; i <= 126; i++ {
		// Each character block is 9 lines high (1 empty line + 8 lines of art)
		start := (i-32)*9 + 1
		asciiMap[rune(i)] = lines[start : start+8]
	}

	// 4. Generate the Art
	result, err := GenerateAscii(input, asciiMap)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 5. Print out the result without an extra newline (the string already formats them)
	fmt.Print(result)
}

// GenerateAscii takes the raw input and ASCII map and builds the final string.
func GenerateAscii(input string, asciiMap map[rune][]string) (string, error) {
	var builder strings.Builder

	// Split by literal "\n" strings passed by bash
	parts := strings.Split(input, "\\n")

	// Edge Case: Check if the string consists ONLY of "\n"
	isAllEmpty := true
	for _, p := range parts {
		if p != "" {
			isAllEmpty = false
			break
		}
	}

	// If it's only "\n"s, print the exact number of newlines requested
	if isAllEmpty {
		for i := 0; i < len(parts)-1; i++ {
			builder.WriteString("\n")
		}
		return builder.String(), nil
	}

	// Normal processing
	for _, part := range parts {
		if part == "" {
			builder.WriteString("\n")
			continue
		}

		// Validate all characters are within printable ASCII scope
		for _, c := range part {
			if c < 32 || c > 126 {
				return "", fmt.Errorf("invalid character %q. only printable ASCII characters allowed", c)
			}
		}

		// Concatenate character blocks line by line (0 to 7)
		for row := 0; row < 8; row++ {
			for _, c := range part {
				builder.WriteString(asciiMap[c][row])
			}
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}
