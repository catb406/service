package logrus

import (
	"bytes"
	"fmt"
	"path/filepath"
)

type GlogLikeFormatter struct {
	Pid string
}

func (f GlogLikeFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%-7s %s ", entry.Level.String(), entry.Time.Format("2006-01-02T15:04:05.000000Z07:00")))
	b.WriteString(f.Pid)
	if entry.HasCaller() {
		fileVal := fmt.Sprintf(" %s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
		if fileVal != "" {
			b.WriteString(fileVal)
		}
	}
	b.WriteString("] ")
	if entry.Message != "" {
		b.WriteString(entry.Message)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}
