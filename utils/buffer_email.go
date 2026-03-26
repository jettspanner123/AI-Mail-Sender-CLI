package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
)

func BuildEmailMessage(config models.SMTPConfig, to, body string, attachmentPaths []string, attachmentDatas [][]byte) ([]byte, error) {
	   var buf bytes.Buffer
	   writer := multipart.NewWriter(&buf)
	   fromHeader := config.From
	   if config.FromName != "" {
			   fromHeader = (&mail.Address{Name: config.FromName, Address: config.From}).String()
	   }

	   fmt.Fprintf(&buf, "From: %s\r\n", fromHeader)
	   fmt.Fprintf(&buf, "To: %s\r\n", to)
	   fmt.Fprintf(&buf, "Subject: %s\r\n", config.Subject)
	   fmt.Fprint(&buf, "MIME-Version: 1.0\r\n")
	   fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%q\r\n", writer.Boundary())
	   fmt.Fprint(&buf, "\r\n")

	   bodyHeaders := textproto.MIMEHeader{}
	   bodyHeaders.Set("Content-Type", "text/plain; charset=\"utf-8\"")
	   bodyHeaders.Set("Content-Transfer-Encoding", "8bit")
	   bodyPart, err := writer.CreatePart(bodyHeaders)
	   if err != nil {
			   return nil, fmt.Errorf("create body part: %w", err)
	   }
	   if _, err := bodyPart.Write([]byte(body)); err != nil {
			   return nil, fmt.Errorf("write body part: %w", err)
	   }

	   for i := range attachmentPaths {
			   if len(attachmentPaths[i]) == 0 || len(attachmentDatas[i]) == 0 {
					   continue // skip empty attachments
			   }
			   attachmentName := filepath.Base(attachmentPaths[i])
			   encodedName := mime.QEncoding.Encode("UTF-8", attachmentName)
			   attachmentHeaders := textproto.MIMEHeader{}
			   attachmentHeaders.Set("Content-Type", "application/pdf; name=\""+encodedName+"\"")
			   attachmentHeaders.Set("Content-Transfer-Encoding", "base64")
			   attachmentHeaders.Set("Content-Disposition", "attachment; filename=\""+encodedName+"\"")
			   attachmentPart, err := writer.CreatePart(attachmentHeaders)
			   if err != nil {
					   return nil, fmt.Errorf("create attachment part: %w", err)
			   }

			   encoded := base64.StdEncoding.EncodeToString(attachmentDatas[i])
			   for j := 0; j < len(encoded); j += 76 {
					   end := j + 76
					   if end > len(encoded) {
							   end = len(encoded)
					   }
					   if _, err := attachmentPart.Write([]byte(encoded[j:end] + "\r\n")); err != nil {
							   return nil, fmt.Errorf("write attachment part: %w", err)
					   }
			   }
	   }

	   if err := writer.Close(); err != nil {
			   return nil, fmt.Errorf("finalize MIME message: %w", err)
	   }

	   return buf.Bytes(), nil
}

func ReadContacts(filePath string) ([]models.Contact, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	nameIdx := FindHeaderIndex(headers, "Name")
	emailIdx := FindHeaderIndex(headers, "Email")
	if nameIdx == -1 || emailIdx == -1 {
		return nil, fmt.Errorf("required headers not found (Name, Email)")
	}

	var contacts []models.Contact
	seen := make(map[string]struct{})

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}

		if nameIdx >= len(record) || emailIdx >= len(record) {
			continue
		}

		email := strings.TrimSpace(record[emailIdx])
		if email == "" {
			continue
		}

		fullName := strings.TrimSpace(record[nameIdx])
		firstName := ExtractFirstName(fullName)
		if firstName == "" {
			continue
		}

		key := strings.ToLower(firstName + "|" + email)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		contacts = append(contacts, models.Contact{
			FullName:  fullName,
			FirstName: firstName,
			Email:     email,
		})
	}

	return contacts, nil
}

func FindHeaderIndex(headers []string, target string) int {
	for i, h := range headers {
		if strings.EqualFold(strings.TrimSpace(h), target) {
			return i
		}
	}
	return -1
}

func ExtractFirstName(fullName string) string {
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
