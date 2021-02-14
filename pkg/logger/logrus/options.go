package logrus

import (
	"github.com/go-impatient/gaia/pkg/logrus/constants"
	"github.com/go-impatient/gaia/pkg/logrus/filehook"
	"github.com/go-impatient/gaia/pkg/logrus/logstashhook"
	"github.com/sirupsen/logrus"
)

type Config struct {
	reportCaller   bool
	level          logrus.Level
	logFormatter   logrus.Formatter
	fields         logrus.Fields
	fileConfig     *filehook.FileConfig
	logstashConfig *logstashhook.LogstashConfig
}

type Option func(*Config)

func NewConfig(opts ...Option) *Config {
	//default
	config := &Config{
		reportCaller:   true,
		level:          logrus.InfoLevel,
		fields:         nil,
		logFormatter:   constants.DefaultTextFormatter,
		fileConfig:     nil,
		logstashConfig: nil,
	}

	for _, opt := range opts {
		opt(config)
	}
	return config
}

func Level(level logrus.Level) Option {
	return func(c *Config) {
		c.level = level
	}
}

func ReportCaller(reportCaller bool) Option {
	return func(c *Config) {
		c.reportCaller = reportCaller
	}
}

func Fields(fields logrus.Fields) Option {
	return func(c *Config) {
		c.fields = fields
	}
}

func LogFormatter(logFormatter logrus.Formatter) Option {
	return func(c *Config) {
		c.logFormatter = logFormatter
	}
}

func FileConfig(fileConfig *filehook.FileConfig) Option {
	return func(c *Config) {
		c.fileConfig = fileConfig
	}
}

func LogstashConfig(logstashConfig *logstashhook.LogstashConfig) Option {
	return func(c *Config) {
		c.logstashConfig = logstashConfig
	}
}
