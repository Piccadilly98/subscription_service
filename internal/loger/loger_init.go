package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var defaultOutPut io.Writer = os.Stdout
var mu sync.RWMutex

func SetDefaultOutput(out io.Writer) error {
	if out == nil {
		return fmt.Errorf("output cannot be nil")
	}

	mu.Lock()
	defer mu.Unlock()
	defaultOutPut = out
	return nil
}

func GetDefaultOutput() io.Writer {
	mu.RLock()
	defer mu.RUnlock()
	return defaultOutPut
}

type LoggerManager struct {
	loggers map[string]*Logger
	mu      sync.RWMutex
}

func NewLoggerManager() *LoggerManager {
	return &LoggerManager{
		loggers: make(map[string]*Logger),
	}
}

func (l *LoggerManager) AddNewLoger(out io.Writer, prefix string) (*Logger, error) {
	if out == nil {
		return nil, fmt.Errorf("output cannot not be nil")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.loggers[prefix]; ok {
		return nil, fmt.Errorf("logger exist, not add new")
	}
	l.loggers[prefix] = newLoger(out, prefix)
	return l.loggers[prefix], nil
}

func (l *LoggerManager) Get(prefix string) *Logger {
	l.mu.RLock()
	loger, ok := l.loggers[prefix]
	l.mu.RUnlock()
	if ok {
		return loger
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	new := newLoger(defaultOutPut, prefix)
	l.loggers[prefix] = new
	return new
}

type Logger struct {
	loger  log.Logger
	prefix string
}

func newLoger(out io.Writer, prefix string) *Logger {
	return &Logger{
		loger:  *log.New(out, fmt.Sprintf("[%s] ", strings.ToUpper(prefix)), log.Ldate|log.Ltime),
		prefix: strings.ToUpper(prefix),
	}
}

func (l *Logger) Printf(level, format string, a ...any) {
	if level != "" {
		fullFormat := level + ": " + format
		l.loger.Printf(fullFormat, a...)
	} else {
		l.loger.Printf(format, a...)
	}
}

func (l *Logger) Println(level string, a ...any) {
	if level != "" {
		message := fmt.Sprint(a...)
		l.loger.Printf("%s: %s\n", level, message)
	} else {
		l.loger.Println(a...)
	}
}

func (l *Logger) ChangeOutput(out io.Writer) {
	l.loger.SetOutput(out)
}
