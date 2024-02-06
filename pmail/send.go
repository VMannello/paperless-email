package pmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"
)

// EmailAccount represents an email account configuration.
type EmailAccount struct {
	Name  string     `yaml:"name"`
	Email string     `yaml:"email"`
	SMTP  SMTPConfig `yaml:"smtp"`
}

// SMTPConfig represents SMTP configuration for an email account.
type SMTPConfig struct {
	Server             string `yaml:"server"`
	Port               int    `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	EnableTLS          bool   `yaml:"tls"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

func SendEmail(cfg Config, tags []string, filePath string) {
	for _, t := range tags {
		msg, ok := cfg.TagMapping[t]
		if !ok {
			log.Printf("[WARN]  tag '%s' not found in configuration, skipping \n", t)
			continue
		}

		var sender EmailAccount
		for _, acc := range cfg.Accounts {
			if acc.Name == msg.From {
				sender = acc
				break
			}
		}

		if sender.Name == "" {
			log.Printf("[WARN]  account [%s] not found in configuration, skipping \n", msg.From)
			continue
		}

		if msg.IncludeAttachment {
			msg.attachments = make(map[string][]byte)
			err := msg.AttachFile(filePath)
			if err != nil {
				log.Printf("[ERROR]  failed to include attachment: %s \n", err)
				return
			}
		}

		if err := sendAttachmentEmail(sender, msg.To, msg.ToBytes()); err != nil {
			log.Printf("[ERROR]  failed to send from [%s] for tag [%s]: %s \n", sender.Name, t, err)
		} else {
			log.Printf("[INFO]  email sent from [%s] for tag [%s] \n", sender.Name, t)
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

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", smtpConfig.Server, smtpConfig.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpConfig.Server)
	if err != nil {
		return fmt.Errorf("failed to build client:  %w", err)
	}
	defer client.Close()

	if tlsConfig != nil {
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start tls: %w", err)
		}
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate:  %w", err)
	}

	if err = client.Mail(acctConfig.Email); err != nil {
		return fmt.Errorf("failed to add sender:  %w", err)
	}

	for _, r := range recipients {
		if err = client.Rcpt(r); err != nil {
			return fmt.Errorf("failed to add recipient:  %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(rawMsg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
