package pmail

import (
	"os"
	"testing"
)

func TestReplaceEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		variables map[string]string // Environment variables to set before the test
		expected  string
	}{
		{
			name:     "empty",
			input:    "This is a test",
			expected: "This is a test",
		},
		{
			name:      "single",
			input:     "{{ENV_VARIABLE}}",
			variables: map[string]string{"ENV_VARIABLE": "TestValue"},
			expected:  "TestValue",
		},
		{
			name:      "multi",
			input:     "Value1: {{VAR1}} Value2: {{VAR2}}",
			variables: map[string]string{"VAR1": "123", "VAR2": "456"},
			expected:  "Value1: 123 Value2: 456",
		},
		{
			name:      "pretty variable",
			input:     "{{ SPACED_VAR }}",
			variables: map[string]string{"SPACED_VAR": "SpaceTest"},
			expected:  "SpaceTest",
		},
		{
			name:      "pretty variable formatted",
			input:     " {{ SPACED_VAR }} ",
			variables: map[string]string{"SPACED_VAR": "SpaceTest"},
			expected:  " SpaceTest ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.variables {
				err := os.Setenv(key, value)
				if err != nil {
					t.Fatal(err)
				}
			}

			result := replaceEnvironmentVariables(tt.input)

			for key := range tt.variables {
				err := os.Unsetenv(key)
				if err != nil {
					t.Fatal(err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}
