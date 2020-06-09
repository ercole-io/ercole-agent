// +build rhel5

package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// BasicLogger struct to compose logger with logrus that satisfy Logger interface
type BasicLogger struct {
	level         Level
	componentName string
	isColored     bool
}

// SetLevel to inner field log
func (l *BasicLogger) SetLevel(level Level) {
	l.level = level
}

// NewLogger return a logrus.Logger initialized with ercole log standard
func NewLogger(componentName string) Logger {
	var newLogger BasicLogger
	newLogger.componentName = componentName
	newLogger.isColored = runtime.GOOS != "windows"
	//newLogg.SetReportCaller(true)
	//newLogg.SetOutput(os.Stdout)

	return &newLogger
}

func (l *BasicLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Debug(msg)
}

func (l *BasicLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Info(msg)
}

func (l *BasicLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Warn(msg)
}

func (l *BasicLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Error(msg)
}

func (l *BasicLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Fatal(msg)
}

func (l *BasicLogger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Panic(msg)
}

func (l *BasicLogger) Debug(args ...interface{}) {
	l.doLog(DebugLevel, args...)
}

func (l *BasicLogger) Info(args ...interface{}) {
	l.doLog(InfoLevel, args...)
}

func (l *BasicLogger) Warn(args ...interface{}) {
	l.doLog(WarnLevel, args...)
}

func (l *BasicLogger) Error(args ...interface{}) {
	l.doLog(ErrorLevel, args...)
}

func (l *BasicLogger) Fatal(args ...interface{}) {
	l.doLog(FatalLevel, args...)

	os.Exit(1)
}

func (l *BasicLogger) Panic(args ...interface{}) {
	l.doLog(PanicLevel, args...)

	s := fmt.Sprintln(args...)
	panic(s)
}

func (l *BasicLogger) doLog(level Level, args ...interface{}) {
	if level < l.level {
		return
	}

	levelColor := getColorByLevel(level))
	levelText := strings.ToUpper(Level.String())[0:4]
	//caller := getCaller(entry)
	caller := ""
	message := strings.TrimSuffix(fmt.Sprintf(args...), "\n")

	var logBuffer bytes.Buffer

	if l.isColored {
		logBuffer.WriteString(fmt.Sprintf("\x1b[%dm", levelColor))
	}

	logBuffer.WriteString(
		fmt.Sprintf("[%s][%s][%s]",
			time.Now().Format("06-01-02 15:04:05"),
			l.componentName
			levelText))

	if l.isColored {
		logBuffer.WriteString("\x1b[0m")
	}

	logBuffer.WriteString(fmt.Sprintf("[%s] %-50s", caller, message))

	for _, k := range getKeysInOrder(entry.Data) {
		logBuffer.WriteString(
			fmt.Sprintf("\x1b[%dm%s\x1b[0m=%v ", levelColor, k, entry.Data[k]))
	}

	buff :=append(logBuffer.Bytes(), '\n'), nil

	log.Println(string(buff))
}
