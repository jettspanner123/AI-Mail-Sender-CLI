package main

import (
	"log/slog"
	"os"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/config"
	"github.com/jettspanner123/AI-Mail-Sender-CLI/utils"
)

func main() {
	logger := slog.New(&utils.PrettyJSONHandler{W: os.Stdout})

	configuration, err := utils.ParseXML()
	if err != nil {
		logger.Error("Failed to parse XML configuration", slog.Any("error", err))
		os.Exit(1)
	}
	contacts, err := utils.ReadContacts(configuration.Assets.Dataset.Path)
	if err != nil {
		logger.Error("Failed to read contacts", slog.Any("error", err), slog.String("path", configuration.Assets.Dataset.Path))
		os.Exit(1)
	}

	smtpConfig, err := config.LoadSMTPConfig()
	if err != nil {
		logger.Error("Invalid SMTP config", slog.Any("error", err))
		os.Exit(1)
	}

	var attachmentsData [][]byte
	var attachmentsPaths []string

	for _, att := range configuration.Assets.Attachments {
		attachmentsPaths = append(attachmentsPaths, att.Path)
	}

	for i, att := range configuration.Assets.Attachments {
		attachmentData, err := os.ReadFile(att.Path)
		if err != nil {
			logger.Error("Failed to read attachment", slog.Int("index", i), slog.String("path", att.Path), slog.Any("error", err))
		}
		attachmentsData = append(attachmentsData, attachmentData)
	}

	logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.Error("Failed to open output.log", slog.Any("error", err))
		os.Exit(1)
	}
	defer logFile.Close()

	logger.Info("Starting email send", slog.Int("contacts", len(contacts)))
	sent, failed := utils.SendBatchEmails(
		logger,
		smtpConfig,
		contacts,
		attachmentsPaths,
		attachmentsData,
		logFile,
	)
	logger.Info("Done sending emails", slog.Int("sent", sent), slog.Int("failed", failed))
}
