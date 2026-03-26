package utils

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
)

func ParseXML() (*models.Configuration, error) {
	file, err := os.Open("/Users/jettspanner123/GoLangDevelopment/ai-mail-sender-cli/assets.xml")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var configuration models.Configuration
	if err := xml.NewDecoder(file).Decode(&configuration); err != nil {
		return nil, err
	}

	configuration.Assets.Subject = strings.TrimSpace(configuration.Assets.Subject)
	configuration.Assets.Body = strings.TrimSpace(configuration.Assets.Body)
	configuration.Assets.Dataset.Name = strings.TrimSpace(configuration.Assets.Dataset.Name)
	configuration.Assets.Dataset.Type = strings.TrimSpace(configuration.Assets.Dataset.Type)
	configuration.Assets.Dataset.Path = strings.TrimSpace(configuration.Assets.Dataset.Path)
	for i := range configuration.Assets.Attachments {
		configuration.Assets.Attachments[i].Name = strings.TrimSpace(configuration.Assets.Attachments[i].Name)
		configuration.Assets.Attachments[i].Type = strings.TrimSpace(configuration.Assets.Attachments[i].Type)
		configuration.Assets.Attachments[i].Path = strings.TrimSpace(configuration.Assets.Attachments[i].Path)
	}

	return &configuration, nil
}

func PrintConfiguration(cfg *models.Configuration) {
	if cfg == nil {
		println("Configuration is nil")
		return
	}
	println("Assets:")
	println("  Dataset:")
	println("    Name:", cfg.Assets.Dataset.Name)
	println("    Type:", cfg.Assets.Dataset.Type)
	println("    Path:", cfg.Assets.Dataset.Path)
	println("  Subject:", cfg.Assets.Subject)
	println("  Body:", cfg.Assets.Body)
	println("  Attachments:")
	for i, att := range cfg.Assets.Attachments {
		println("    Attachment", i+1, ":")
		println("      Name:", att.Name)
		println("      Type:", att.Type)
		println("      Path:", att.Path)
	}
}
