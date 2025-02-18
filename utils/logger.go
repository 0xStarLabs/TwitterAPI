package utils

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

type Logger struct {
	level LogLevel
}

// LogLevel defines the verbosity of logging
type LogLevel int

const (
	LogLevelNone    LogLevel = iota // No logging
	LogLevelError                   // Only errors
	LogLevelWarning                 // Errors and warnings
	LogLevelSuccess                 // Errors, warnings, and success messages
	LogLevelInfo                    // Normal operational logs
	LogLevelDebug                   // Detailed debug information
)

func NewLogger(level LogLevel) Logger {
	return Logger{level: level}
}

func (l Logger) Debug(format string, a ...any) error {
	if l.level >= LogLevelDebug {
		color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=blue;op=bold> | DEBUG   | </><fg=white>-</> %s\n",
			time.Now().Format("15:04:05.000"),
			fmt.Sprintf(format, a...))
	}
	return fmt.Errorf(format, a...)
}

func (l Logger) Info(format string, a ...any) error {
	if l.level >= LogLevelInfo {
		color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=white;op=bold> | INFO    | </><fg=white>-</> %s\n",
			time.Now().Format("15:04:05.000"),
			fmt.Sprintf(format, a...))
	}
	return fmt.Errorf(format, a...)
}

func (l Logger) Error(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	if l.level >= LogLevelError {
		color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=red;op=bold> | ERROR   | </><fg=white>-</> %s\n",
			time.Now().Format("15:04:05.000"),
			err.Error())
	}
	return err
}

func (l Logger) Success(format string, a ...any) error {
	if l.level >= LogLevelSuccess {
		color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=green;op=bold> | SUCCESS | </><fg=white>-</> %s\n",
			time.Now().Format("15:04:05.000"),
			fmt.Sprintf(format, a...))
	}
	return fmt.Errorf(format, a...)
}

func (l Logger) Warning(format string, a ...any) error {
	if l.level >= LogLevelWarning {
		color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=yellow;op=bold> | WARNING | </><fg=white>-</> %s\n",
			time.Now().Format("15:04:05.000"),
			fmt.Sprintf(format, a...))
	}
	return fmt.Errorf(format, a...)
}
