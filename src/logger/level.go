package logger

type Level int8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	PERF
)

func (v Level) Enable(lvl Level) bool {
	return v <= lvl
}
