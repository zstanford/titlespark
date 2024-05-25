package util

import (
	"unicode"
)

func RemoveSpecialCharacters(input string) string {
	var result []rune
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func IsOnlySpaces(input string) bool {
	for _, char := range input {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}
