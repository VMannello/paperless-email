package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/vmannello/paperless-email/internal/config"
)

// SendEmail sends an email with an attachment.
func SendEmail(cfg config.Config, tags []string, filePath string) {
	// Process each tag mapping
	for _, tag := range tags {
		mapCfg, ok := cfg.TagMapping.Config[tag]
		if !ok {
			log.Printf("[WARN]  tag '%s' not found in configuration, skipping", tag)
			continue
		}

		var selectedAccount config.EmailAccount
		for _, acc := range cfg.Accounts {
			if acc.Name == mapCfg.Account {
				selectedAccount = acc
				break
			}
		}

		// ensure the account was found
		if selectedAccount.Name == "" {
			log.Printf("[WARN]  account [%s] not found in configuration, skipping", mapCfg.Account)
			continue
		}

		if err := sendAttachmentEmail(selectedAccount, cfg.Message, mapCfg.Recipients, filePath); err != nil {
			log.Printf("[ERROR]  failed to send from [%s] to [%s] for tag [%s]: %s", selectedAccount.Name, mapCfg.Recipients, tag, err)
		} else {
			log.Printf("[INFO]  email sent from [%s] to [%s] for tag [%s]", selectedAccount.Name, mapCfg.Recipients, tag)
		}
	}
}

func sendAttachmentEmail(acctConfig config.EmailAccount, msgConfig config.Message, recipients []string, filePath string) error {
	smtpConfig := acctConfig.SMTP
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Server)

	var tlsConfig *tls.Config
	if smtpConfig.EnableTLS {
		tlsConfig = &tls.Config{
			ServerName:         smtpConfig.Server,
			InsecureSkipVerify: smtpConfig.InsecureSkipVerify,
		}
	}

	conn, err := smtp.Dial(smtpConfig.Server + ":" + fmt.Sprint(smtpConfig.Port))
	if err != nil {
		return fmt.Errorf("error connecting to SMTP server: %w", err)
	}
	defer conn.Close()

	if tlsConfig != nil {
		if err := conn.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("error starting tls: %w", err)
		}
	}

	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("error authenticating:  %w", err)
	}

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s",
		strings.Join(recipients, ","), msgConfig.Subject, msgConfig.Body)

	if msgConfig.IncludeAttachment {
		attachmentName := filepath.Base(filePath)
		message += "\r\n--BOUNDARY\r\n" +
			"Content-Type: application/octet-stream\r\n" +
			"Content-Disposition: attachment; filename=\"" + attachmentName + "\"\r\n\r\n"

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// Append the file content to the message
		message += string(fileContent)

		// Add closing boundary
		message += "\r\n--BOUNDARY--"
	}

	// Send the email
	err = smtp.SendMail(smtpConfig.Server+":"+fmt.Sprint(smtpConfig.Port), auth, acctConfig.Email, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
