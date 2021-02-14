package logstashhook

import (
	"fmt"
	"github.com/go-impatient/gaia/pkg/logrus/constants"
	"github.com/go-impatient/gaia/pkg/tcp"
	"github.com/sirupsen/logrus"
	"io"
	"path"
	"runtime"
	"strings"
	"sync"
)

type LogstashConfig struct {
	logFormatter logrus.Formatter
	fields       logrus.Fields
	address      string
	sourceField  string
	skip         int
	open         bool
}

type Option func(*LogstashConfig)

func NewLogstashConfig(address string, opts ...Option) *LogstashConfig {
	//default
	config := &LogstashConfig{
		open:         true,
		logFormatter: constants.DefaultJSONFormatter,
		fields:       nil,
		address:      address,
		sourceField:  "source",
		skip:         10,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (lc *LogstashConfig) IsOpen() bool {
	return lc.open
}

func Open(open bool) Option {
	return func(c *LogstashConfig) {
		c.open = open
	}
}

func LogFormatter(logFormatter logrus.Formatter) Option {
	return func(c *LogstashConfig) {
		c.logFormatter = logFormatter
	}
}

func Fields(fields logrus.Fields) Option {
	return func(c *LogstashConfig) {
		c.fields = fields
	}
}

func SourceField(sourceField string) Option {
	return func(c *LogstashConfig) {
		c.sourceField = sourceField
	}
}

func Skip(skip int) Option {
	return func(c *LogstashConfig) {
		c.skip = skip
	}
}

func (lc *LogstashConfig) validate() error {
	if lc.address == "" {
		return fmt.Errorf("you must configure address in [Conf.LogstashConf]; %#v", lc)
	}
	return nil
}

func (lc *LogstashConfig) CreateLogstashHook() (logrus.Hook, error) {
	err := lc.validate()
	if err != nil {
		return nil, err
	}
	conn, err := tcp.Dial("tcp", lc.address)
	if err != nil {
		return nil, fmt.Errorf("net.Dial(tcp, %s); Error : %v", lc.address, err)
	}
	return Hook{
		sourceField: lc.sourceField,
		skip:        lc.skip,
		writer:      conn,
		formatter: LogstashFormatter{
			Formatter: lc.logFormatter,
			Fields:    lc.fields,
		},
	}, nil
}

type Hook struct {
	sourceField string
	skip        int
	writer      io.Writer
	formatter   logrus.Formatter
}

func (h Hook) Fire(e *logrus.Entry) error {
	e.Data[h.sourceField] = findCaller(h.skip)
	//e.Data["level"] = strings.ToUpper(e.Level.String())
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)
	if err != nil {
		return err
	}
	return nil
}

func findCaller(skip int) string {
	result := ""
	for i := 1; i <= skip; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			funcName := runtime.FuncForPC(pc).Name()
			result = fmt.Sprintf("%s:%s:%d", path.Base(file), path.Base(funcName), line)
			if !strings.Contains(funcName, "logrus") && !strings.Contains(funcName, "logrus") {
				break
			}
			//fmt.Println(fmt.Sprintf("line:%d; pc : %s; file : %s; funcName : %s; path.Base(funcName) : %s; path.Base(file) : %s", line, pc, file, funcName, path.Base(funcName), path.Base(file)))
		}
	}
	return result
}

func (h Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type LogstashFormatter struct {
	logrus.Formatter
	logrus.Fields
}

func (f LogstashFormatter) Format(e *logrus.Entry) ([]byte, error) {
	ne := copyEntry(e, f.Fields)

	dataBytes, err := f.Formatter.Format(ne)

	releaseEntry(ne)
	return dataBytes, err
}

func copyEntry(e *logrus.Entry, fields logrus.Fields) *logrus.Entry {
	ne := entryPool.Get().(*logrus.Entry)
	ne.Message = e.Message
	ne.Level = e.Level
	ne.Time = e.Time
	ne.Data = logrus.Fields{}
	for k, v := range fields {
		ne.Data[k] = v
	}
	for k, v := range e.Data {
		ne.Data[k] = v
	}
	return ne
}

var entryPool = sync.Pool{
	New: func() interface{} {
		return &logrus.Entry{}
	},
}

func releaseEntry(e *logrus.Entry) {
	entryPool.Put(e)
}
