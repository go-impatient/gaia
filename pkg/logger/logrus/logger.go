package logrus

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/go-impatient/gaia/pkg/logger/logrus/fieldhook"
)

// Represents loggo logger and fields to support structured logging.
type Logger struct {
	entry  *logrus.Entry
	config *Config
}

// Initiate logger.
func New(config *Config) *Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	e := logrus.NewEntry(logger)
	l := Logger{
		entry:  e,
		config: config,
	}
	l.EnableTextLogging()
	return &l
}

func (l *Logger) validate() error {
	return nil
}

func (l *Logger) Init() {
	log.Println(fmt.Sprintf("InitialLogConfig : %#v", l))
	err := l.validate()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(l.config.level)
	logrus.SetReportCaller(l.config.reportCaller)
	logrus.SetFormatter(l.config.logFormatter)
	if l.config.fields != nil && len(l.config.fields) > 0 {
		logrus.AddHook(&fieldhook.DefaultFieldHook{Fields: l.config.fields})
	}
	err = l.initialFileConfig()
	if err != nil {
		panic(err)
	}
	err = l.initialLogstashConf()
	if err != nil {
		panic(err)
	}
}

func (l *Logger) initialFileConfig() error {
	if l.config.fileConfig == nil {
		return nil
	}
	if !l.config.fileConfig.IsOpen() {
		return nil
	}
	infoHook, err := l.config.fileConfig.CreateFileHook("info", []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel})
	if err != nil {
		return fmt.Errorf("create info hook error; Error : %v", err)
	}
	logrus.AddHook(infoHook)
	errorHook, err := l.config.fileConfig.CreateFileHook("error", []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel})
	if err != nil {
		return fmt.Errorf("create file hook error; Error : %v", err)
	}
	logrus.AddHook(errorHook)
	return nil
}

func (l *Logger) initialLogstashConf() error {
	if l.config.logstashConfig == nil {
		return nil
	}
	if !l.config.logstashConfig.IsOpen() {
		return nil
	}
	hook, err := l.config.logstashConfig.CreateLogstashHook()
	if err != nil {
		return fmt.Errorf("create logstash hook error; Error : %v", err)
	}
	logrus.AddHook(hook)
	return nil
}

// Log correct file name and line number from where Logger call was invoked.
func prettyfier(r *runtime.Frame) (string, string) {
	return "", ""
}

func (l *Logger) EnableJSONLogging() {
	l.entry.Logger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: prettyfier,
	})
}

func (l *Logger) EnableTextLogging() {
	l.entry.Logger.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier:       prettyfier,
		DisableLevelTruncation: true,
	})
}
