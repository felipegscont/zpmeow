package logger

import (
	waLog "go.mau.fi/whatsmeow/util/log"
)

// waLogAdapter adapts our Logger interface to implement waLog.Logger
type waLogAdapter struct {
	logger Logger
}

// NewWALogAdapter creates a new adapter that implements waLog.Logger interface
func NewWALogAdapter(logger Logger) waLog.Logger {
	return &waLogAdapter{
		logger: logger,
	}
}

// Warnf implements waLog.Logger interface
func (w *waLogAdapter) Warnf(msg string, args ...interface{}) {
	w.logger.Warnf(msg, args...)
}

// Errorf implements waLog.Logger interface
func (w *waLogAdapter) Errorf(msg string, args ...interface{}) {
	w.logger.Errorf(msg, args...)
}

// Infof implements waLog.Logger interface
func (w *waLogAdapter) Infof(msg string, args ...interface{}) {
	w.logger.Infof(msg, args...)
}

// Debugf implements waLog.Logger interface
func (w *waLogAdapter) Debugf(msg string, args ...interface{}) {
	w.logger.Debugf(msg, args...)
}

// Sub implements waLog.Logger interface
func (w *waLogAdapter) Sub(module string) waLog.Logger {
	return &waLogAdapter{
		logger: w.logger.Sub(module),
	}
}

// GetWALogger returns a waLog.Logger that uses our logger system
func GetWALogger(module string) waLog.Logger {
	return NewWALogAdapter(GetLogger().Sub(module))
}
