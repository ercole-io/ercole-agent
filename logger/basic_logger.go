// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
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

// NewBasicLogger return a BasicLogger initialized with ercole log standard
func NewBasicLogger(componentName string) Logger {
	var newLogger BasicLogger
	newLogger.componentName = componentName
	newLogger.isColored = runtime.GOOS != "windows"

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
	if level > l.level {
		return
	}

	levelColor := getColorByLevel(level)
	levelText := strings.ToUpper(Level.String(level))[0:4]
	message := strings.TrimSuffix(fmt.Sprint(args...), "\n")

	var buffer bytes.Buffer

	if l.isColored {
		buffer.WriteString(fmt.Sprintf("\x1b[%dm", levelColor))
	}

	buffer.WriteString(
		fmt.Sprintf("[%s][%s][%s]",
			time.Now().Format("06-01-02 15:04:05"),
			l.componentName,
			levelText))

	if l.isColored {
		buffer.WriteString("\x1b[0m")
	}

	buffer.WriteString(fmt.Sprintf("%-50s", message))

	log.Println(buffer.String())
}
