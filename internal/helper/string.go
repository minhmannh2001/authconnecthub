package helper

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ToLowerFirstChar(str string) string {
	if str == "" {
		return str
	}

	return strings.ToLower(string(str[0])) + str[1:]
}

func FormatToastMessage(sentence string) string {
	// Split the sentence into words using hyphen (-) as delimiter
	words := strings.Split(sentence, "-")

	// Handle empty sentence case
	if len(words) == 1 && words[0] == "" {
		return ""
	}

	// Uppercase the first letter of the first word
	words[0] = strings.ToUpper(string(words[0][0])) + words[0][1:]

	for i, word := range words {
		if i > 0 {
			if strings.HasSuffix(words[i-1], ".") {
				words[i] = strings.ToUpper(string(word[0])) + word[1:]
			}
		}
	}

	// Join the words with spaces and return the formatted sentence
	return strings.Join(words, " ")
}

func MapToJSONString(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling map to JSON: %w", err)
	}
	return string(jsonData), nil
}
