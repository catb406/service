package logrus

import (
	"bytes"
	"fmt"
	"path/filepath"
)

type OnlyMessageFormatter struct{}

func (f OnlyMessageFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	if entry.HasCaller() {
		fileVal := fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
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
