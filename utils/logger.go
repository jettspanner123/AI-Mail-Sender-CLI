package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func WriteResult(logFile *os.File, line string) {
	fmt.Println(line)
	if _, err := logFile.WriteString(line + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write log line: %v\n", err)
	}
}

type PrettyJSONHandler struct {
	W io.Writer
}

func (h *PrettyJSONHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *PrettyJSONHandler) Handle(_ context.Context, r slog.Record) error {
	m := make(map[string]interface{})
	m["time"] = r.Time
	m["level"] = r.Level.String()
	m["msg"] = r.Message
	r.Attrs(func(a slog.Attr) bool {
		m[a.Key] = a.Value.Any()
		return true
	})
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = h.W.Write(append(b, '\n'))
	return err
}
func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler       { return h }
