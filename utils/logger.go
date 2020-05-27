package utils

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewLogger return a logrus.Logger initialized with ercole log standard
func NewLogger(componentName string) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&ercoleFormatter{
		ComponentName: componentName[0:4],
		isColored:     runtime.GOOS != "windows",
	})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	return logger
}

// ercoleFormatter custom formatter for ercole that formats logs into text
type ercoleFormatter struct {
	ComponentName string
	isColored     bool
}

// Format renders a single log entry
func (f *ercoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry)
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	caller := getCaller(entry)
	message := strings.TrimSuffix(entry.Message, "\n")

	var logBuffer bytes.Buffer

	if f.isColored {
		logBuffer.WriteString(fmt.Sprintf("\x1b[%dm", levelColor))
	}

	logBuffer.WriteString(
		fmt.Sprintf("[%s][%s][%s]",
			entry.Time.Format("06-01-02 15:04:05"),
			f.ComponentName,
			levelText))

	if f.isColored {
		logBuffer.WriteString("\x1b[0m")
	}

	logBuffer.WriteString(fmt.Sprintf("[%s] %-50s", caller, message))

	for _, k := range getKeysInOrder(entry.Data) {
		logBuffer.WriteString(
			fmt.Sprintf("\x1b[%dm%s\x1b[0m=%v ", levelColor, k, entry.Data[k]))
	}

	return append(logBuffer.Bytes(), '\n'), nil
}

func getColorByLevel(entry *logrus.Entry) int {
	const gray = 37
	const yellow = 33
	const red = 31
	const blue = 36

	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return gray
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	default:
		return blue
	}
}

func getCaller(entry *logrus.Entry) string {
	if !entry.HasCaller() {
		return ""
	}

	caller := entry.Caller.File
	if strings.Contains(caller, "ercole-agent/") {
		caller = caller[strings.Index(caller, "ercole-agent/")+len("ercole-agent/"):]
	}

	return fmt.Sprintf("%s:%d", caller, entry.Caller.Line)
}

func getKeysInOrder(entryData logrus.Fields) []string {
	manuallyOrderedKeys := []string{"endpoint", "statusCode"}

	for i := 0; i < len(manuallyOrderedKeys); i++ {
		if _, ok := entryData[manuallyOrderedKeys[i]]; !ok {
			manuallyOrderedKeys = remove(manuallyOrderedKeys, i)
			i--
		}
	}

	var entryDataKeys []string
	for k := range entryData {
		if !contains(manuallyOrderedKeys, k) {
			entryDataKeys = append(entryDataKeys, k)
		}
	}
	sort.Strings(entryDataKeys)

	return append(manuallyOrderedKeys, entryDataKeys...)
}

// contains return true if a contains x, otherwise false.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// remove return slice without element at position i, mantaining order
func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}
