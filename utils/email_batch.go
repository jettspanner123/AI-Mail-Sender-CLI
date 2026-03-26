package utils

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jettspanner123/AI-Mail-Sender-CLI/constants"
	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
)

// SendBatchEmails handles concurrent batch sending and logging of emails.
func SendBatchEmails(
	logger *slog.Logger,
	smtpConfig models.SMTPConfig,
	contacts []models.Contact,
	attachmentsPaths []string,
	attachmentsData [][]byte,
	logFile *os.File,
) (sent, failed int) {
	const concurrentSends = 10
	successSincePause := 0

	for i := 0; i < len(contacts); i += concurrentSends {
		end := i + concurrentSends
		if end > len(contacts) {
			end = len(contacts)
		}

		batch := contacts[i:end]
		results := make(chan models.SendResult, len(batch))
		var wg sync.WaitGroup

		for _, contact := range batch {
			c := contact
			wg.Add(1)
			go func() {
				defer wg.Done()

				personalizedBody := strings.ReplaceAll(constants.EMAIL_CONTENT, "{{First Name}}", c.FirstName)
				msg, err := BuildEmailMessage(
					smtpConfig,
					c.Email,
					personalizedBody,
					attachmentsPaths,
					attachmentsData,
				)
				if err != nil {
					logger.Error("Failed to build email message", slog.String("email", c.Email), slog.Any("error", err))
					results <- models.SendResult{
						Line:    fmt.Sprintf("%s::%s::UnSuccessful ❌::TimeTaken 0s", c.FullName, c.Email),
						Success: false,
					}
					return
				}

				duration, err := SendEmailWithProgress(smtpConfig, c.Email, msg)
				if err != nil {
					logger.Error("Failed to send email", slog.String("email", c.Email), slog.Any("error", err), slog.Duration("duration", duration))
					results <- models.SendResult{
						Line:    fmt.Sprintf("%s::%s::UnSuccessful ❌::TimeTaken %s", c.FullName, c.Email, duration.Round(time.Second)),
						Success: false,
					}
					return
				}

				logger.Info("Email sent successfully", slog.String("email", c.Email), slog.Duration("duration", duration))
				results <- models.SendResult{
					Line:    fmt.Sprintf("%s::%s::Successful ✅::TimeTaken %s", c.FullName, c.Email, duration.Round(time.Second)),
					Success: true,
				}
			}()
		}

		wg.Wait()
		close(results)

		for result := range results {
			if result.Success {
				sent++
				successSincePause++
			} else {
				failed++
			}
			logger.Info("Result", slog.String("result", result.Line))
			WriteResult(logFile, result.Line)
		}

		for successSincePause >= 10 {
			logger.Info("Reached 10 successful emails. Waiting 1 minute before continuing...")
			time.Sleep(1 * time.Minute)
			successSincePause -= 10
		}
	}

	return sent, failed
}

