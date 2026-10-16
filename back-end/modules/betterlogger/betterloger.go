package betterlogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Level string

const (
	LevelInfo    Level = "INFO"
	LevelSuccess Level = "SUCCESS"
	LevelWarn    Level = "WARN"
	LevelError   Level = "ERROR"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m" // Info
	colorGreen  = "\033[32m" // Success
	colorYellow = "\033[33m" // Warn
	colorRed    = "\033[31m" // Error
)

type Logger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	file          *os.File
	mu            sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

// Init initializes the singleton logger with an optional file destination.
func Init(filePath string) error {
	var initErr error
	once.Do(func() {
		var writers []io.Writer
		writers = append(writers, os.Stdout)

		var logFile *os.File
		if filePath != "" {
			var err error
			logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				initErr = err
				return
			}
			writers = append(writers, logFile)
		}

		// Console gets standard logger with flags (date, time, short file)
		consoleWriter := io.MultiWriter(os.Stdout)

		// If a file is provided, we set up a separate file writer to avoid ANSI color clutter
		var fileWriter io.Writer
		if logFile != nil {
			fileWriter = logFile
		}

		instance = &Logger{
			consoleLogger: log.New(consoleWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
			fileLogger:    log.New(fileWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
			file:          logFile,
		}
	})
	return initErr
}

// GetInstance returns the singleton logger instance, initializing with stdout only if not already initialized.
func GetInstance() *Logger {
	if instance == nil {
		_ = Init("")
	}
	return instance
}

// logMessage handles the formatting, colorization, and multi-destination routing.
func (l *Logger) logMessage(level Level, color string, format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, v...)

	// Formatted strings for console (with color) and file (plain text)
	consoleOutput := fmt.Sprintf("%s[%s]%s %s", color, level, colorReset, msg)
	fileOutput := fmt.Sprintf("[%s] %s", level, msg)

	// Output to console
	_ = l.consoleLogger.Output(3, consoleOutput)

	// Output to file if configured
	if l.fileLogger != nil {
		_ = l.fileLogger.Output(3, fileOutput)
	}
}

// Public logging methods
func (l *Logger) Info(format string, v ...any) {
	l.logMessage(LevelInfo, colorBlue, format, v...)
}

func (l *Logger) Success(format string, v ...any) {
	l.logMessage(LevelSuccess, colorGreen, format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.logMessage(LevelWarn, colorYellow, format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.logMessage(LevelError, colorRed, format, v...)
}

// Close safely closes the log file if open.
func Close() {
	if instance != nil && instance.file != nil {
		instance.mu.Lock()
		defer instance.mu.Unlock()
		_ = instance.file.Close()
	}
}
