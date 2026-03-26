package utils

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
)

func SendEmail(config models.SMTPConfig, to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	var auth smtp.Auth
	if config.Username != "" && config.Password != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	return smtp.SendMail(addr, auth, config.From, []string{to}, msg)
}

func SendEmailWithProgress(config models.SMTPConfig, email string, msg []byte) (time.Duration, error) {
	start := time.Now()

	err := SendEmail(config, email, msg)

	return time.Since(start), err
}
