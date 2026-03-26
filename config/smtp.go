package config

import (
	"errors"
	"os"
	"strings"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/constants"
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
		Subject:  constants.EMAIL_SUBJECT,
	}

	if config.Port == "" {
		config.Port = constants.SMTP_DEFAULT_MAIL_PORT
	}

	if config.Host == "" {
		config.Host = constants.SMTP_DEFAULT_HOST
	}

	if config.From == "" {
		return models.SMTPConfig{}, errors.New("SMTP_FROM is empty! Exiting...")
	}

	return config, nil
}
