package utils

// Package utils provides utility functions and modules like reusable logger
// helper functions for ErrorWrappers, ToPtrs etc.

// The logger.go file implements a wrapper around logrus that provides
// a singleton logger instance with various configuration options.
// It supports logging to console and file simultaneously, different log levels,
// and both text and JSON formatting.

import (
	"io"
	"os"
	"sync"
	"time"

	logrus "github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	files []*os.File
	mu    sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

func WithLogLevel(level string) func(*Logger) {
	return func(log *Logger) {
		logLevel, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.Fatal("Failed to parse log level: ", err)
		}
		log.SetLevel(logLevel)
	}
}

func WithLogTextFormatter(format *logrus.TextFormatter) func(*Logger) {
	return func(log *Logger) {
		log.SetFormatter(format)
	}
}

func WithLogJSONFormatter(format *logrus.JSONFormatter) func(*Logger) {
	return func(log *Logger) {
		log.SetFormatter(format)
	}
}

func WithFileOutput(filePath string) func(*Logger) {
	return func(log *Logger) {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			logrus.Fatal("Failed to open log file: ", err)
		}
		log.mu.Lock()
		log.files = append(log.files, file)
		log.mu.Unlock()
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	}
}

func GetLogger(opts ...func(*Logger)) *Logger {
	once.Do(func() {
		instance = &Logger{
			Logger: logrus.New(),
			files:  make([]*os.File, 0),
		}

		instance.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			PadLevelText:    true,
		})

		for _, opt := range opts {
			opt(instance)
		}
	})
	return instance
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lastErr error
	for _, file := range l.files {
		if err := file.Sync(); err != nil {
			lastErr = err
		}
		if err := file.Close(); err != nil {
			lastErr = err
		}
	}
	l.files = nil
	return lastErr
}