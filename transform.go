package main

import (
	"strings"
	"unicode"
)

func MapsToEnumData(maps []MapParseResult) []EnumData {
	var ed []EnumData

	for _, mapdata := range maps {
		ed = append(ed, mapToEnumData(mapdata))
	}
	return ed
}

func mapToEnumData(pm MapParseResult) EnumData {
	var ed EnumData
	ed.SourceName = pm.Name
	ed.Type = toPascalCase(strings.TrimSuffix(pm.Name, "Values"))
	ed.Value = strings.ToLower(ed.Type) + "value"
	ed.Description = ed.Type + " simulates an enum with available values."
	ed.Items = map[string]string{}

	for value, displayname := range pm.Values {
		ed.Items[value] = toPascalCase(displayname)
	}
	return ed
}

func toPascalCase(s string) string {
	// Handle empty string early to avoid unnecessary processing
	if s == "" {
		return s
	}

	// First, split the string into words, handling multiple delimiter types
	words := splitIntoWords(s)

	// Build the final string using strings.Builder for better performance
	var result strings.Builder
	result.Grow(len(s)) // Pre-allocate capacity to avoid reallocations

	for _, word := range words {
		if word == "" {
			continue
		}

		// Special handling for acronyms - preserve their case
		if isAcronym(word) {
			result.WriteString(strings.ToUpper(word))
			continue
		}

		// Capitalize first letter, lowercase rest
		result.WriteString(capitalize(word))
	}

	return result.String()
}

// splitIntoWords breaks a string into words, handling both camelCase and
// delimiter-separated strings. It maintains acronym boundaries.
func splitIntoWords(s string) []string {
	var words []string
	var currentWord strings.Builder
	var prevChar rune

	for i, char := range s {
		switch {
		// Handle delimiters (space, hyphen, underscore)
		case unicode.IsSpace(char) || char == '-' || char == '_':
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}

		// Handle camelCase boundaries
		case i > 0 && unicode.IsUpper(char) && !unicode.IsUpper(prevChar):
			// Add current word if exists and start new word
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			currentWord.WriteRune(char)

		// Handle normal characters
		default:
			currentWord.WriteRune(char)
		}
		prevChar = char
	}

	// Add the last word if it exists
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// isAcronym checks if a word is likely an acronym by verifying if
// it's all uppercase letters.
func isAcronym(s string) bool {
	for _, char := range s {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return len(s) > 1 // Require at least 2 characters to be considered an acronym
}

// capitalize converts the first character of a string to uppercase
// and the rest to lowercase.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}
