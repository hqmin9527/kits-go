package logger

import "testing"

func TestCaller(t *testing.T) {
	Debug("normal")
	NewLogger("", "").Info("NewLogger")
	Prefix("").Debug("LogPrefix")
	GetLogger().Info("GetLogger")
	GetLoggerWithSkip(-1).Info("GetLoggerWithSkip(-1)")
}
