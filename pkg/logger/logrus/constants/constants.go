package constants

import (
	"time"

	"github.com/sirupsen/logrus"
)

var DefaultTextFormatter = &logrus.TextFormatter{
	DisableColors:    false,
	FullTimestamp:    true,
	TimestampFormat:  string("2006-01-02 15:04:05"),
	QuoteEmptyFields: true,
	FieldMap: *&logrus.FieldMap{
		logrus.FieldKeyMsg:  "message",
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyFunc: "method",
		logrus.FieldKeyFile: "source",
		//logrus.FieldKeyLevel: "@level",
	},
}

var DefaultJSONFormatter = &logrus.JSONFormatter{
	TimestampFormat: time.RFC3339Nano,
	FieldMap: *&logrus.FieldMap{
		logrus.FieldKeyMsg:  "message",
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyFunc: "method",
		logrus.FieldKeyFile: "source",
		//logrus.FieldKeyLevel: "@level",
	},
}
