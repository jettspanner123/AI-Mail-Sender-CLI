package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
)

func LoadSMTPConfig() (models.SMTPConfig, error) {
	config := models.SMTPConfig{
		Host:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port:     strings.TrimSpace(os.Getenv("SMTP_PORT")),
		Username: strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
		From:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
		Subject:  strings.TrimSpace(os.Getenv("SMTP_SUBJECT")),
	}

	if config.Port == "" {
		config.Port = "587"
	}
	if config.Subject == "" {
		config.Subject = "Hunar.ai | Conversational AI Recruiters"
	}

	if config.Host == "" || config.From == "" {
		return models.SMTPConfig{}, fmt.Errorf("SMTP_HOST and SMTP_FROM are required")
	}

	return config, nil
}
