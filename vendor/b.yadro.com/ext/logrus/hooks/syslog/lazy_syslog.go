package syslog

import (
	"errors"
	"flag"
	"fmt"
	"log/syslog"
	"os"
	"path/filepath"
	"sync"

	"b.yadro.com/ext/logrus"
)

var (
	programname string
)

func init() {
	flag.StringVar(&programname, "programname", filepath.Base(os.Args[0]), "used to tag syslog messages")
}

// LazySyslogHook to send logs via syslog.
type LazySyslogHook struct {
	mu            sync.Mutex
	writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
	priority      syslog.Priority
	formatter     logrus.Formatter
}

// NewLazySyslogHook creates a hook to be added to an instance of logger.
// Initialization of syslog writer happens only when Fire() is called.
// This is a workaround for initialization of tag by a command-line option flag
func NewLazySyslogHook(network, raddr string, priority syslog.Priority, fmtr logrus.Formatter) *LazySyslogHook {
	return &LazySyslogHook{
		writer:        nil,
		SyslogNetwork: network,
		SyslogRaddr:   raddr,
		priority:      priority,
		formatter:     fmtr,
	}
}

func (hook *LazySyslogHook) init() (err error) {
	if hook.writer == nil {
		hook.mu.Lock()
		defer hook.mu.Unlock()
		if hook.writer == nil {
			if !flag.Parsed() {
				return errors.New("initializing syslog hook before flag.Parse")
			}
			hook.writer, err = syslog.Dial(hook.SyslogNetwork, hook.SyslogRaddr, hook.priority, programname)
		}
	}
	return
}

func (hook *LazySyslogHook) Fire(entry *logrus.Entry) error {
	if err := hook.init(); err != nil {
		return err
	}
	lineBytes, err := hook.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to format entry, %v", err)
		return err
	}
	line := string(lineBytes)
	switch entry.Level {
	case logrus.PanicLevel:
		return hook.writer.Crit(line)
	case logrus.FatalLevel:
		return hook.writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.writer.Err(line)
	case logrus.WarnLevel:
		return hook.writer.Warning(line)
	case logrus.InfoLevel:
		return hook.writer.Info(line)
	case logrus.DebugLevel, logrus.TraceLevel:
		return hook.writer.Debug(line)
	default:
		return nil
	}
}

func (hook *LazySyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
