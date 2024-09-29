// Log configuration

package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"kevinsheeran/walrus-backend/config"
	"os"
	"path"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	log := config.Config.Logger

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		fmt.Fprintf(b, "%s [%s] [%dm][%s] [%s] ", log.Prefix, timestamp, levelColor, funcVal, fileVal)
	} else {
		fmt.Fprintf(b, "%s [%s] ", timestamp, levelColor)
	}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	mLogger := logrus.New()
	mLogger.SetOutput(os.Stdout)
	mLogger.SetReportCaller(config.Config.Logger.ShowLineNumber)
	mLogger.SetFormatter(&LogFormatter{})
	level, err := logrus.ParseLevel(config.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	mLogger.SetLevel(level)
	InitDefaultLogger()
	return mLogger
}

func InitDefaultLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(config.Config.Logger.ShowLineNumber)
	logrus.SetFormatter(&LogFormatter{})
	level, err := logrus.ParseLevel(config.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}
