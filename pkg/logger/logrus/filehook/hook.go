package filehook

import (
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/go-impatient/gaia/pkg/logrus/constants"
)

type FileConfig struct {
	logFormatter     logrus.Formatter
	filePath         string
	withMaxAge       time.Duration
	withRotationTime time.Duration
	open             bool
}

func NewFileConfig(filePath string, opts ...Option) *FileConfig {
	//default
	config := &FileConfig{
		open:             true,
		logFormatter:     constants.DefaultTextFormatter,
		filePath:         filePath,
		withMaxAge:       time.Duration(876000) * time.Hour,
		withRotationTime: time.Duration(24) * time.Hour,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (lc *FileConfig) IsOpen() bool {
	return lc.open
}

type Option func(*FileConfig)

func Open(open bool) Option {
	return func(c *FileConfig) {
		c.open = open
	}
}

func LogFormatter(logFormatter logrus.Formatter) Option {
	return func(c *FileConfig) {
		c.logFormatter = logFormatter
	}
}

func WithMaxAge(withMaxAge time.Duration) Option {
	return func(c *FileConfig) {
		c.withMaxAge = withMaxAge
	}
}

func WithRotationTime(withRotationTime time.Duration) Option {
	return func(c *FileConfig) {
		c.withRotationTime = withRotationTime
	}
}

func createFolder(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (fc *FileConfig) validate() error {
	if fc.filePath == "" {
		return fmt.Errorf("you must configure filePath;FileConfig : %#v", fc)
	}
	if !exists(fc.filePath) {
		err := createFolder(fc.filePath)
		if err != nil {
			return fmt.Errorf("create folder error; Error : %v", err)
		}
	}
	return nil
}

func (fc *FileConfig) CreateFileHook(fileName string, writerLevels []logrus.Level) (*lfshook.LfsHook, error) {
	err := fc.validate()
	if err != nil {
		return nil, err
	}
	hookWrite, err := rotatelogs.New(
		fc.filePath+fileName+".log.%Y%m%d",
		rotatelogs.WithLinkName(fc.filePath+fileName+".log"),
		//rotatelogs.WithLinkName(lc.filePath+lc.module+"-info.log"),
		rotatelogs.WithMaxAge(fc.withMaxAge),
		rotatelogs.WithRotationTime(fc.withRotationTime),
	)
	if err != nil {
		return nil, err
	}
	if writerLevels == nil || len(writerLevels) == 0 {
		return nil, fmt.Errorf("writer levels can not be empty")
	}
	writerMap := make(lfshook.WriterMap)
	for _, writerLevel := range writerLevels {
		writerMap[writerLevel] = hookWrite
	}
	return lfshook.NewHook(writerMap, fc.logFormatter), nil
}
