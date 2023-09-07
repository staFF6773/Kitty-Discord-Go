package log

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

// Logging Level
type Level int

const (
	// LogDebug represents debug messages.
	LogDebug Level = iota

	// LogInfo represents informational messages.
	LogInfo

	// LogError represents errors.
	LogError

	// LogFatal represents fatal errors.
	LogFatal
)

var (
	// LogLevelDisplayNames gives the display name to use for our log levels.
	LogLevelDisplayNames = map[Level]string{
		LogDebug: "debug",
		LogInfo:  "info",
		LogError: "error",
		LogFatal: "fatal",
	}

	LogLevelNames = map[string]Level{
		"debug":  LogDebug,
		"info":   LogInfo,
		"error":  LogError,
		"errors": LogError,
		"fatal":  LogFatal,
	}
)

func (logger *Logger) Debug(logType string, messageParts ...string) {
	logger.Log(LogDebug, logType, messageParts...)
}

func (logger *Logger) Info(logType string, messageParts ...string) {
	logger.Log(LogInfo, logType, messageParts...)
}

func (logger *Logger) Error(logType string, messageParts ...string) {
	logger.Log(LogError, logType, messageParts...)
}

func (logger *Logger) Fatal(logType string, messageParts ...string) {
	logger.Log(LogFatal, logType, messageParts...)
	os.Exit(1)
}

type fileMethod struct {
	Enabled  bool
	Filename string
	File     *os.File
	Writer   *bufio.Writer
}

type Logger struct {
	stdoutWriteLock *sync.Mutex
	fileWriteLock   *sync.Mutex
	MethodSTDOUT    bool
	MethodSTDERR    bool
	MethodFile      fileMethod
	Level           Level
}

func NewLogger(method []string, file string, level string) (*Logger, error) {
	logger := Logger{
		MethodFile: fileMethod{
			Filename: file,
		},
		Level: LogLevelNames[level],
	}

	for _, p := range method {
		switch p {
		//?case syslog
		case "stdout":
			logger.MethodSTDOUT = true
			logger.stdoutWriteLock = &sync.Mutex{}

		case "stderr":
			logger.MethodSTDERR = true
			logger.stdoutWriteLock = &sync.Mutex{}

		case "file":
			logger.MethodFile.Enabled = true
			logger.fileWriteLock = &sync.Mutex{}

			if logFile, fileErr := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm|os.ModeAppend); fileErr != nil {
				return nil, fileErr
			} else {
				logger.MethodFile.File = logFile
				logger.MethodFile.Writer = bufio.NewWriter(logFile)
			}
		}
	}
	return &logger, nil
}

func (logger *Logger) Close() error {
	if logger.MethodFile.Enabled {
		flushErr := logger.MethodFile.Writer.Flush()
		closeErr := logger.MethodFile.File.Close()
		if flushErr != nil {
			return flushErr
		}
		return closeErr
	}
	return nil
}

func (logger *Logger) Log(level Level, logType string, messageParts ...string) {
	// no logging enabled
	if !(logger.MethodSTDOUT || logger.MethodSTDERR || logger.MethodFile.Enabled) {
		return
	}

	// check if we log on the right level
	if level < logger.Level {
		return
	}

	// assemble our log line
	var rawBuf bytes.Buffer
	fmt.Fprintf(&rawBuf, "[%s]-[%s]-[%s] ", time.Now().Format("2006-01-02T15:04:05.000Z"), LogLevelDisplayNames[level], logType)
	for i, p := range messageParts {
		rawBuf.WriteString(p)

		if i != len(messageParts)-1 {
			rawBuf.WriteString(" : ")
		}
	}
	rawBuf.WriteRune('\n')

	// output
	if logger.MethodSTDOUT {
		logger.stdoutWriteLock.Lock()
		os.Stdout.Write(rawBuf.Bytes())
		logger.stdoutWriteLock.Unlock()
	}

	if logger.MethodSTDERR {
		logger.stdoutWriteLock.Lock()
		os.Stderr.Write(rawBuf.Bytes())
		logger.stdoutWriteLock.Unlock()
	}

	if logger.MethodFile.Enabled {
		logger.fileWriteLock.Lock()
		logger.MethodFile.Writer.Write(rawBuf.Bytes())
		logger.MethodFile.Writer.Flush()
		logger.fileWriteLock.Unlock()
	}
}
