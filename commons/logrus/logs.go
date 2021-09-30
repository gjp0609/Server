package logrus

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var log = logrus.New()

// MyFormatter 自定义 formatter
type MyFormatter struct {
}

// Format implement the Formatter interface
func (myFormatter *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	// entry.Message 就是需要打印的日志
	b.WriteString(fmt.Sprintf("%s %s [%20s:%-4d] %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"),
		getLevelName(entry.Level),
		filepath.Base(entry.Caller.File),
		entry.Caller.Line,
		entry.Message))
	return b.Bytes(), nil
}

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
	//log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&MyFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func GetLogger() *logrus.Logger {
	return log
}

func getLevelName(level logrus.Level) string {
	var levelName string
	switch level {
	case logrus.TraceLevel:
	case logrus.DebugLevel:
		levelName = "DEBUG"
	case logrus.InfoLevel:
		levelName = "INFO"
	case logrus.WarnLevel:
		levelName = "WARN"
	case logrus.ErrorLevel:
		levelName = "ERROR"
	case logrus.PanicLevel:
		levelName = "PANIC"
	case logrus.FatalLevel:
		levelName = "FATAL"
	}
	return fmt.Sprintf("%5s", levelName)
}
