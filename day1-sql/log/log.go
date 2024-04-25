package sqllog

import (
	"io"
	"log"
	"os"
	"sync"
)

// loggers
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// SetLevel sets the log level for the loggers.
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 判断当前level能使什么级别的日志输出
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	// init loggers
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	// disable logs
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
