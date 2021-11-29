package logger

import (
	"fmt"
	"log/syslog"

	log "b.yadro.com/ext/logrus"
	logrus_syslog "b.yadro.com/ext/logrus/hooks/syslog"
)

type LogLevel string

const (
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warning"
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
	TraceLevel LogLevel = "trace"

	loggerStackFrameOffset = 1
)

func init() {
	hook := logrus_syslog.NewLazySyslogHook("", "", syslog.LOG_LOCAL0, log.OnlyMessageFormatter{})
	log.AddHook(hook)
	log.SetReportCaller(true)
	log.SetFormatter(log.GlogLikeFormatter{})
}

func format(args []interface{}) string {
	var resultString string = ""
	if len(args) > 0 {
		var format string
		format, ok := args[0].(string)
		if !ok {
			str, ok := args[0].(fmt.Stringer)
			if ok {
				format = str.String()
			}
		}
		newArgs := args[1:]
		resultString = fmt.Sprintf(format, newArgs...)
	}
	return resultString
}

func SetLogLevel(l LogLevel) error {
	return SetLogLevelString(string(l))
}

func SetLogLevelString(l string) error {
	level, err := log.ParseLevel(l)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}

func IsDebug() bool {
	return log.IsLevelEnabled(log.DebugLevel)
}

func IsTrace() bool {
	return log.IsLevelEnabled(log.TraceLevel)
}

// Panic logs a message and calls panic()
func Panic(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Panic(format(args))
}

// Fatal logs a message and calls os.Exit(1)
func Fatal(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Fatal(format(args))
}

func Error(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Error(format(args))
}

func Warn(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Warning(format(args))
}

func Info(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Info(format(args))
}

func Debug(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Debug(format(args))
}

func Trace(args ...interface{}) {
	log.WithSkipCallers(loggerStackFrameOffset).Trace(format(args))
}

func PanicDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Panic(format(args))
}

func FatalDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Fatal(format(args))
}

func ErrorDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Error(format(args))
}

func WarningDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Warning(format(args))
}

func InfoDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Info(format(args))
}

func DebugDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Debug(format(args))
}

func TraceDepth(depth int, args ...interface{}) {
	log.WithSkipCallers(depth + loggerStackFrameOffset).Trace(format(args))
}
