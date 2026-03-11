package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	config "github.com/jettspanner123/AI-Mail-Sender-CLI/config"
	constants "github.com/jettspanner123/AI-Mail-Sender-CLI/constants"
	"github.com/jettspanner123/AI-Mail-Sender-CLI/models"
	utils "github.com/jettspanner123/AI-Mail-Sender-CLI/utils"
)

func main() {
	filePath := "assets/dataset.csv"
	attachmentPath := "assets/Hunar Conversational AI Agents_Self Serve_V1.pdf"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if len(os.Args) > 2 {
		attachmentPath = os.Args[2]
	}

	contacts, err := utils.ReadContacts(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read contacts: %v\n", err)
		os.Exit(1)
	}

	config, err := config.LoadSMTPConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid SMTP config: %v\n", err)
		os.Exit(1)
	}

	attachmentData, err := os.ReadFile(attachmentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read attachment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sending %d emails...\n", len(contacts))
	sent := 0
	failed := 0
	successSincePause := 0

	logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output.log: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	const concurrentSends = 10
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
				msg, err := utils.BuildEmailMessage(config, c.Email, personalizedBody, attachmentPath, attachmentData)
				if err != nil {
					results <- models.SendResult{
						Line:    fmt.Sprintf("%s::%s::UnSuccessful ❌::TimeTaken 0s", c.FullName, c.Email),
						Success: false,
					}
					return
				}

				duration, err := utils.SendEmailWithProgress(config, c.Email, msg)
				if err != nil {
					results <- models.SendResult{
						Line:    fmt.Sprintf("%s::%s::UnSuccessful ❌::TimeTaken %s", c.FullName, c.Email, duration.Round(time.Second)),
						Success: false,
					}
					return
				}

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
			utils.WriteResult(logFile, result.Line)
		}

		for successSincePause >= 10 {
			fmt.Println("Reached 10 successful emails. Waiting 1 minute before continuing...")
			time.Sleep(1 * time.Minute)
			successSincePause -= 10
		}
	}

	fmt.Printf("Done. Sent: %d, Failed: %d\n", sent, failed)
}
