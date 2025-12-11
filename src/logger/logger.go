package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	logLevel      = DEBUG
	defaultLogger = newLoggerCore("", "", 5)
)

// 创建新的Logger，不能保证并发安全（相同的输出可能会打乱日志顺序）
func NewLogger(logDir, logFile string) *Logger {
	return newLoggerCore(logDir, logFile, 4)
}

// 可以直接用来打印
func GetLogger() *Logger {
	return GetLoggerWithSkip(-1)
}

// 一般用于底层库封装，以便打印真实调用的文件行号
func GetLoggerWithSkip(skip int) *Logger {
	return defaultLogger.addCallerSkip(skip)
}

func SetLogLevel(lvl Level) {
	logLevel = lvl
}
func SetLogFile(logDir, logFile string) {
	defaultLogger.SetLogFile(logDir, logFile)
}
func End() {
	defaultLogger.End()
}

func Debug(format string, args ...any) {
	defaultLogger.Debug(format, args...)
}
func Info(format string, args ...any) {
	defaultLogger.Info(format, args...)
}
func Warn(format string, args ...any) {
	defaultLogger.Warn(format, args...)
}
func Error(format string, args ...any) {
	defaultLogger.Error(format, args...)
}
func Perf(format string, args ...any) {
	defaultLogger.Perf(format, args...)
}

type Logger struct {
	core       **log.Logger // 输出核心
	callerSkip int          // 调用栈
	// 输出到文件时属性
	byFile   bool     // 输出到文件
	fPath    []byte   // 文件地址
	fPtr     *os.File // 文件句柄
	dayBegin int64    // 按日切分，当前文件日期开始时间戳
}

func newLoggerCore(logDir, logFile string, skip int) *Logger {
	l := new(Logger)
	l.callerSkip = skip
	l.SetLogFile(logDir, logFile)
	return l
}

func (l *Logger) SetLogFile(logDir, logFile string) {
	var w io.Writer
	switch logFile {
	case "tty", "stdout", "":
		w = os.Stdout
	default:
		if logDir != "" {
			_ = os.MkdirAll(logDir, 0755)
		}
		l.byFile = true
		logFile = strings.TrimSuffix(logFile, ".log")
		l.fPath = []byte(filepath.Join(logDir, logFile))
		l.updateFilePtr(time.Now())
		w = l.fPtr
	}
	lo := log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	// lo := log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	if l.core == nil {
		l.core = &lo
	} else {
		*l.core = lo
	}
}
func (l *Logger) End() {
}

func (l *Logger) addCallerSkip(skip int) *Logger {
	return &Logger{
		core:       l.core,
		callerSkip: l.callerSkip + skip,
	}
}

func (l *Logger) Debug(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logCore(DEBUG, l.logFnc, msg)
}
func (l *Logger) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logCore(INFO, l.logFnc, msg)
}
func (l *Logger) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logCore(WARN, l.logFnc, msg)
}
func (l *Logger) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logCore(ERROR, l.logFnc, msg)
}
func (l *Logger) Perf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	logCore(PERF, l.logFnc, msg)
}

func (l *Logger) logFnc(msg string) {
	// 更新输出文件
	if l.byFile && l.updateFilePtr(time.Now()) {
		(*l.core).SetOutput(l.fPtr)
	}
	_ = (*l.core).Output(l.callerSkip, msg)
}

func (l *Logger) updateFilePtr(now time.Time) bool {
	if now.Unix()-l.dayBegin < 24*3600 {
		return false
	}
	for now.Unix()-l.dayBegin >= 24*3600 {
		l.dayBegin += 24 * 3600
		if l.fPtr != nil {
			_ = l.fPtr.Close()
		}
		l.fPtr = nil
	}
	if l.fPtr == nil {
		var err error
		name := now.AppendFormat(l.fPath, "2006-01-02.log")
		l.fPtr, err = os.OpenFile(string(name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("open log file failed: %s", err)
		}
	}
	return true
}

func logCore(lvl Level, lf logFnc, msg string) {
	if !logLevel.Enable(lvl) {
		return
	}
	msg = logLevelPrefix(lvl) + msg
	lf(msg)
}

func logLevelPrefix(lvl Level) string {
	switch lvl {
	case DEBUG:
		return "[DEBUG] "
	case INFO:
		return "[INFO] "
	case WARN:
		return "[WARN] "
	case ERROR:
		return "[ERROR] "
	case PERF:
		return "[PERF] "
	default:
		return "[LEVEL(" + strconv.Itoa(int(lvl)) + ")]"
	}
}

type logFnc func(v string)

func DayBegin() time.Time {
	now := time.Now()
	_, offset := now.Zone()
	s := now.Unix()
	s = s - (s+int64(offset))%(24*60*60)
	return time.Unix(s, 0)
}
