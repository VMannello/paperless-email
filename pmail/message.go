package pmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Message struct {
	From              string   `yaml:"from"`
	To                []string `yaml:"to"`
	CC                []string `yaml:"cc"`
	BCC               []string `yaml:"bcc"`
	Subject           string   `yaml:"subject"`
	Body              string   `yaml:"body"`
	IncludeAttachment bool     `yaml:"include_attachment"`

	attachments map[string][]byte
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	m.Body = replaceEnvironmentVariables(m.Body)

	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

var curlyRe = regexp.MustCompile(`\{\{([^\}]+)\}\}`)

func replaceEnvironmentVariables(input string) string {
	return curlyRe.ReplaceAllStringFunc(input, func(match string) string {
		return os.Getenv(strings.TrimSpace(strings.Trim(match, "{}")))
	})
}
