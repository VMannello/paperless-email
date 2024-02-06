package pmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
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
