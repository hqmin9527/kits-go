package logger

type Prefix string

var prefixLogger = GetLoggerWithSkip(1)

type logFunc func(format string, v ...any)

func (l Prefix) Debug(format string, v ...any) {
	l.log(prefixLogger.Debug, format, v)
}
func (l Prefix) Info(format string, v ...any) {
	l.log(prefixLogger.Info, format, v)
}
func (l Prefix) Warn(format string, v ...any) {
	l.log(prefixLogger.Warn, format, v)
}
func (l Prefix) Error(format string, v ...any) {
	l.log(prefixLogger.Error, format, v)
}
func (l Prefix) Perf(format string, v ...any) {
	l.log(prefixLogger.Perf, format, v)
}

func (l Prefix) SetLogFile(logDir, logFile string) {
	prefixLogger.SetLogFile(logDir, logFile)
}
func (l Prefix) End() {
	prefixLogger.End()
}

func (l Prefix) log(lf logFunc, format string, v []any) {
	if l != "" {
		format = string(l) + " " + format
	}
	lf(format, v...)
}

func (l Prefix) Append(ap string) Prefix {
	return l + Prefix(" ["+ap+"]")
}
