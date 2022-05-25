// LIMITAION: can only prints interface{}
package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LogTypeInfo  = "INFO"
	LogTypeWarn  = "WARN"
	LogTypeError = "ERROR"
	LogTypeFatal = "FATAL"
)

type l interface {
	Error()
	Info()
	Fatal()
}

type ll struct {
	Error *log.Logger
	Info  *log.Logger
	Fatal *log.Logger
}

func Newl(w io.Writer) *ll {
	return &ll{
		Error: log.New(w, "Error", log.Lshortfile|log.LstdFlags),
		Info:  log.New(w, "Info", log.Lshortfile|log.LstdFlags),
		Fatal: log.New(w, "Fatal", log.Lshortfile|log.LstdFlags),
	}
}

type Logger struct {
	timeFormat   string
	minSkipTrace int
	minLogLevel  int // TODO
}

func New(timeFormat string, minSkipTrace int) *Logger {
	return &Logger{
		timeFormat:   timeFormat,
		minSkipTrace: minSkipTrace,
	}
}

func (l Logger) formatMessage(severity string, data interface{}) string {
	now := time.Now().Format(l.timeFormat)
	_, file, line, _ := runtime.Caller(l.minSkipTrace)
	file = filepath.Base(file)
	return fmt.Sprintf("%s: %s: %s:%d %s",
		severity,
		now,
		file,
		line,
		data,
	)
}

func (l Logger) print(out io.Writer, severity string, data interface{}, a ...interface{}) {
	message := l.formatMessage(severity, data)
	if len(a) > 0 {
		fmt.Fprintf(out, message+"\n", a...)
	} else {
		fmt.Fprintln(out, message)
	}

	if severity == LogTypeFatal {
		os.Exit(1)
	}
}

func (l Logger) Error(message interface{}, a ...interface{}) {
	l.print(os.Stderr, LogTypeError, message, a...)
}

func (l Logger) Warn(message interface{}, a ...interface{}) {
	l.print(os.Stdout, LogTypeWarn, message, a...)
}
func (l Logger) Info(message interface{}, a ...interface{}) {
	l.print(os.Stdout, LogTypeInfo, message, a...)
}

func (l Logger) Fatal(message interface{}, a ...interface{}) {
	l.print(os.Stderr, LogTypeFatal, message, a...)
}

func (l Logger) ServerError(w http.ResponseWriter, message interface{}, a ...interface{}) {
	l.print(os.Stderr, LogTypeError, message, a...)
	w.WriteHeader(http.StatusInternalServerError)
}
