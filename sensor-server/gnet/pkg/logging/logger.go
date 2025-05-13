// Copyright (c) 2020 The Gnet Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logging provides logging functionality for gnet applications,
// it sets up a default logger (powered by go.uber.org/zap)
// that is about to be used by your gnet application.
// You're allowed to replace the default logger with your customized logger by
// implementing Logger and assign it to the functional option via gnet.WithLogger,
// and then passing it to gnet.Run or gnet.Rotate.
//
// The environment variable `GNET_LOGGING_LEVEL` determines which zap logger level will be applied for logging.
// The environment variable `GNET_LOGGING_FILE` is set to a local file path when you want to print logs into local file.
// Alternatives of logging level (the variable of logging level ought to be integer):
/*
const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)
*/
package logging

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jurgen-kluft/go-home/sensor-server/gnet/pkg/logging/logf"
)

// Logger is used for logging formatted messages.
type Logger interface {
	// Debugf logs messages at DEBUG level.
	Debugf(format string, args ...any)
	// Infof logs messages at INFO level.
	Infof(format string, args ...any)
	// Warnf logs messages at WARN level.
	Warnf(format string, args ...any)
	// Errorf logs messages at ERROR level.
	Errorf(format string, args ...any)
	// Fatalf logs messages at FATAL level.
	Fatalf(format string, args ...any)
}

// Flusher is the callback function which flushes any buffered log entries to the underlying writer.
// It is usually called before the gnet process exits.
type Flusher = func() error

var (
	mu                  sync.RWMutex
	defaultLogger       Logger
	defaultLoggingLevel Level
	defaultFlusher      Flusher
)

// Level is the alias of logf.Level.
type Level = logf.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = logf.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = logf.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = logf.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = logf.ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = logf.FatalLevel
)

func init() {
	lvl := os.Getenv("GNET_LOGGING_LEVEL")
	if len(lvl) > 0 {
		loggingLevel, err := strconv.ParseInt(lvl, 10, 8)
		if err != nil {
			panic("invalid GNET_LOGGING_LEVEL, " + err.Error())
		}
		defaultLoggingLevel = Level(loggingLevel)
	}

	logger := logf.New(logf.Opts{
		EnableColor:          true,
		Level:                logf.DebugLevel,
		CallerSkipFrameCount: 3,
		EnableCaller:         true,
		TimestampFormat:      time.RFC3339Nano,
		DefaultFields:        []interface{}{"scope", "example"},
	})

	defaultLogger = logger
}

func CreateLoggerAsLocalFile(logLevel Level) (Logger, Flusher, error) {
	// Implementation of CreateLoggerAsLocalFile
	// This function should create a logger that writes to a local file
	// and return the logger and its flusher.
	// For now, we will just return nil for the logger and flusher.
	return defaultLogger, func() error {
		return nil
	}, nil
}

// GetDefaultLogger returns the default logger.
func GetDefaultLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLogger
}

// GetDefaultFlusher returns the default flusher.
func GetDefaultFlusher() Flusher {
	mu.RLock()
	defer mu.RUnlock()
	return defaultFlusher
}

// SetDefaultLoggerAndFlusher sets the default logger and its flusher.
func SetDefaultLoggerAndFlusher(logger Logger, flusher Flusher) {
	mu.Lock()
	defaultLogger, defaultFlusher = logger, flusher
	mu.Unlock()
}

// LogLevel tells what the default logging level is.
func LogLevel() string {
	return strings.ToUpper(defaultLoggingLevel.String())
}

// Cleanup does something windup for logger, like closing, flushing, etc.
func Cleanup() {
	mu.RLock()
	if defaultFlusher != nil {
		_ = defaultFlusher()
	}
	mu.RUnlock()
}

// Error prints err if it's not nil.
func Error(err error) {
	if err != nil {
		mu.RLock()
		defaultLogger.Errorf("error occurs during runtime, %v", err)
		mu.RUnlock()
	}
}

// Debugf logs messages at DEBUG level.
func Debugf(format string, args ...any) {
	mu.RLock()
	defaultLogger.Debugf(format, args...)
	mu.RUnlock()
}

// Infof logs messages at INFO level.
func Infof(format string, args ...any) {
	mu.RLock()
	defaultLogger.Infof(format, args...)
	mu.RUnlock()
}

// Warnf logs messages at WARN level.
func Warnf(format string, args ...any) {
	mu.RLock()
	defaultLogger.Warnf(format, args...)
	mu.RUnlock()
}

// Errorf logs messages at ERROR level.
func Errorf(format string, args ...any) {
	mu.RLock()
	defaultLogger.Errorf(format, args...)
	mu.RUnlock()
}

// Fatalf logs messages at FATAL level.
func Fatalf(format string, args ...any) {
	mu.RLock()
	defaultLogger.Fatalf(format, args...)
	mu.RUnlock()
}
