package utils

import (
	"fmt"
	"os"
)

func WriteResult(logFile *os.File, line string) {
	fmt.Println(line)
	if _, err := logFile.WriteString(line + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write log line: %v\n", err)
	}
}
