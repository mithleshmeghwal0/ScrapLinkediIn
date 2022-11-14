package log

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	logCallers = true
)

func New() *logrus.Entry {
	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{
		DataKey:     "data",
		PrettyPrint: true,
	})

	l.SetReportCaller(logCallers)

	l.Level = logrus.TraceLevel

	return logrus.NewEntry(l)
}

func NewWithFile(f io.Writer) *logrus.Entry {
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{
		DataKey:         "data",
		TimestampFormat: time.RFC3339Nano,
	}

	l.SetReportCaller(logCallers)

	l.Level = logrus.TraceLevel

	l.SetOutput(f)

	return logrus.NewEntry(l)
}
