package helper_test

import (
	"fmt"
	"testing"

	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestToLowerFirstChar(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "Hello", expected: "hello"},
		{input: "WORLD", expected: "wORLD"},
		{input: "123ABC", expected: "123ABC"},
		{input: "MixedCase123", expected: "mixedCase123"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := helper.ToLowerFirstChar(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormatToastMessage(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "success", expected: "Success"},
		{input: "this is a formatted message", expected: "This is a formatted message"},
		{input: "this is a formatted message. another sentence", expected: "This is a formatted message. another sentence"},
		{input: "user-created", expected: "User created"},
		{input: "user-created.-another-message", expected: "User created. Another message"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := helper.FormatToastMessage(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapToJSONString_EmptyMap(t *testing.T) {
	data := map[string]interface{}{}
	jsonString, err := helper.MapToJSONString(data)
	assert.Nil(t, err)
	assert.Equal(t, `{}`, jsonString)
}

func TestMapToJSONString_SimpleMap(t *testing.T) {
	data := map[string]interface{}{"name": "John", "age": 30}
	expectedJSON := "{\"age\":30,\"name\":\"John\"}"
	jsonString, err := helper.MapToJSONString(data)
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, jsonString)
}

func TestMapToJSONString_MarshalError(t *testing.T) {
	invalidData := map[string]interface{}{"error": func() {}} // Function as value will cause marshal error
	_, err := helper.MapToJSONString(invalidData)
	assert.NotNil(t, err)
	expectedError := "json: unsupported type: func()"
	assert.Equal(t, err.Error(), fmt.Sprintf("error marshalling map to JSON: %s", expectedError))
}

func TestRandStringBytes_Length(t *testing.T) {
	// Test cases with different desired lengths
	testCases := []struct {
		n    int
		want int
	}{
		{n: 5, want: 5},
		{n: 10, want: 10},
		{n: 20, want: 20},
	}

	for _, tc := range testCases {
		result := helper.RandStringBytes(tc.n)
		if len(result) != tc.want {
			t.Errorf("Expected string of length %d, got %d", tc.want, len(result))
		}
	}
}

func TestRandStringBytes_Content(t *testing.T) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Test if the generated string contains characters from the letterBytes slice
	result := helper.RandStringBytes(10)
	for _, char := range result {
		found := false
		for _, letter := range letterBytes {
			if char == letter {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Generated string contains character not present in letterBytes")
		}
	}
}
