package pmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

func SendEmail(cfg Config, tags []string, filePath string) {
	for _, tag := range tags {
		mapCfg, ok := cfg.TagMapping.Config[tag]
		if !ok {
			log.Printf("[WARN]  tag '%s' not found in configuration, skipping", tag)
			continue
		}

		var selectedAccount EmailAccount
		for _, acc := range cfg.Accounts {
			if acc.Name == mapCfg.Account {
				selectedAccount = acc
				break
			}
		}

		if selectedAccount.Name == "" {
			log.Printf("[WARN]  account [%s] not found in configuration, skipping", mapCfg.Account)
			continue
		}

		if !cfg.Message.IncludeAttachment {
			filePath = ""
		}

		rawMsg, err := buildMessage(mapCfg.Recipients, cfg.Message.Subject, cfg.Message.Body, filePath)
		if err != nil {
			log.Printf("[ERROR]  error building message: %s", err)
			break
		}

		if err := sendAttachmentEmail(selectedAccount, mapCfg.Recipients, rawMsg); err != nil {
			log.Printf("[ERROR]  failed to send from [%s] to [%s] for tag [%s]: %s", selectedAccount.Name, mapCfg.Recipients, tag, err)
		} else {
			log.Printf("[INFO]  email sent from [%s] to [%s] for tag [%s]", selectedAccount.Name, mapCfg.Recipients, tag)
		}
	}
}

func sendAttachmentEmail(acctConfig EmailAccount, recipients []string, rawMsg []byte) error {
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
	
	err = smtp.SendMail(smtpConfig.Server+":"+fmt.Sprint(smtpConfig.Port), auth, acctConfig.Email, recipients, rawMsg)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func buildMessage(recipients []string, subject, body, attachmentPath string) ([]byte, error) {
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s",
		strings.Join(recipients, ","), subject, body)

	if attachmentPath != "" {
		attachmentName := filepath.Base(attachmentPath)
		message += "\r\n--BOUNDARY\r\n" +
			"Content-Type: application/octet-stream\r\n" +
			"Content-Disposition: attachment; filename=\"" + attachmentName + "\"\r\n\r\n"

		fileContent, err := os.ReadFile(attachmentPath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		message += string(fileContent)
		message += "\r\n--BOUNDARY--"
	}

	return []byte(message), nil
}