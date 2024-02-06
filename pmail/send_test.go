package pmail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	tests := []struct {
		name            string
		recipients      []string
		subject, body   string
		attachmentPath  string
		expectedMessage string
		expectError     bool
	}{
		{
			name:            "No Attachment",
			recipients:      []string{"recipient@example.com"},
			subject:         "Test Subject",
			body:            "Test Body",
			attachmentPath:  "",
			expectedMessage: "To: recipient@example.com\r\nSubject: Test Subject\r\n\r\nTest Body",
			expectError:     false,
		},
		{
			name:            "With Attachment",
			recipients:      []string{"recipient@example.com"},
			subject:         "Test Subject",
			body:            "Test Body",
			attachmentPath:  createTempAttachment(),
			expectedMessage: "To: recipient@example.com\r\nSubject: Test Subject\r\n\r\nTest Body\r\n--BOUNDARY\r\nContent-Type: application/octet-stream\r\nContent-Disposition: attachment; filename=\"attachment.txt\"\r\n\r\nAttachment Content\r\n--BOUNDARY--",
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := buildMessage(tt.recipients, tt.subject, tt.body, tt.attachmentPath)

			if tt.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.expectedMessage, string(message), "Unexpected message content")
			}
		})
	}
}

func createTempAttachment() string {
	content := "Attachment Content"
	tmpfile, err := os.CreateTemp("", "attachment.*.txt")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		panic(err)
	}

	return tmpfile.Name()
}
