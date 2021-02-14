package logrus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/go-impatient/gaia/pkg/logger/logrus/filehook"
	"github.com/go-impatient/gaia/pkg/logger/logrus/logstashhook"
)

var defaultFields = logrus.Fields{
	"project":  "test",
	"module":   "test-module",
	"@version": "1",
	"type":     "project-log",
}

func TestDefaultLog(t *testing.T) {
	var config = NewConfig(
		Fields(defaultFields),
	)
	log := New(config)

	log.Init()
	logrus.Debugf("test debug log[%s]", "This is debug log")
	logrus.Infof("test info log[%s]", "This is info log")
	logrus.Errorf("test error log[%s]", "This is error log")
	logrus.Warningf("test warn log[%s]", "This is warn log")

}

func TestFileLog(t *testing.T) {
	var fileConfig = filehook.NewFileConfig("/Users/moocss/work/GoProjects/gaia/pkg/logrus/logs/")
	var config = NewConfig(
		Fields(defaultFields),
		FileConfig(fileConfig),
	)
	log := New(config)
	log.Init()
	logrus.Debugf("test debug log[%s]", "This is debug log")
	logrus.Infof("test info log[%s]", "This is info log")
	logrus.Errorf("test error log[%s]", "This is error log")
	logrus.Warningf("test warn log[%s]", "This is warn log")
}

func TestLogstashLog(t *testing.T) {
	var logstashConfig = logstashhook.NewLogstashConfig("localhost:5000", logstashhook.Fields(defaultFields))
	var config = NewConfig(
		Fields(defaultFields),
		LogstashConfig(logstashConfig),
	)
	log := New(config)
	log.Init()
	logrus.Debugf("test debug log[%s]", "This is debug log")
	logrus.Infof("test info log[%s]", "111111111111")
	logrus.Errorf("test error log[%s]", "22222222222")
	logrus.Warningf("test warn log[%s]", "This is warn log")
}

//func TestInputLogstash(t *testing.T) {
//	var logstashConfig = logstashhook.NewLogstashConfig("localhost:5000", logstashhook.Fields(defaultFields))
//	var config = NewLogrusConfig(
//		Fields(defaultFields),
//		LogstashConfig(logstashConfig),
//	)
//	config.Initial()
//	var i int64
//	for {
//		logrus.Infof("test info[%d] %v", i, time.Now().Format("2006-01-02 15:04:05"))
//		i++
//		time.Sleep(5 * time.Second)
//	}
//}

// Customize Formatter
func TestCustomizeFormatter(t *testing.T) {
	var customizeFormatter = &CustomizeFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		Levels:          []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"},
		Project:         "unknown-project",
		Module:          "unknown-module",
		Version:         "unknown-version",
		Debug:           false,
	}

	var config = NewConfig(Fields(defaultFields), LogFormatter(customizeFormatter))
	log := New(config)

	log.Init()

	logrus.Infof("test info log[%s]", "This is info log")
}

// Test concurrent modifications to fields.
func TestConcurrentMods(t *testing.T) {
	var count = 15000
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	fileName := filepath.Join(dir, fmt.Sprintf("tc%d", time.Now().Unix()))
	fmt.Println(fileName)
	//defer os.Remove(fileName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0777)
	require.NoError(t, err)
	defer f.Close()
	var config = NewConfig(
		Fields(defaultFields),
	)
	log := New(config)

	log.Init()
	log.entry.Logger.SetOutput(f)
	if level, err := logrus.ParseLevel("debug"); err == nil {
		log.entry.Logger.SetLevel(level)
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			iStr := fmt.Sprintf("%d", i)
			log.entry.WithError(errors.New("No alert matched in alert config.")).WithField("key"+iStr, "value"+iStr).Debugf("Run %d", i)
		}(i)
	}
	wg.Wait()
	_, _ = f.Seek(0, 0)
	output, err := ioutil.ReadAll(f)
	require.NoError(t, err)
	for i := 0; i < count; i++ {
		require.Contains(t, string(output), fmt.Sprintf("level=%s msg=\"Run %d\" error=\"No alert matched in alert config.\" key%d=value%d", "debug", i, i, i))
	}
}

// Test concurrent modifications to fields with JSON formatter.
func TestJSONConcurrentMods(t *testing.T) {
	var count = 15000
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	fileName := filepath.Join(dir, fmt.Sprintf("tc%d", time.Now().Unix()))
	defer os.Remove(fileName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0777)
	require.NoError(t, err)
	defer f.Close()
	var config = NewConfig(
		Fields(defaultFields),
	)
	log := New(config)

	log.Init()
	log.entry.Logger.SetOutput(f)
	if level, err := logrus.ParseLevel("debug"); err == nil {
		log.entry.Logger.SetLevel(level)
	}
	log.EnableJSONLogging()
	wg := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			iStr := fmt.Sprintf("%d", i)
			log.entry.WithError(errors.New("No alert matched in alert config.")).WithField("key"+iStr, "value"+iStr).Debugf("Run %d", i)
		}(i)
	}
	wg.Wait()
	_, _ = f.Seek(0, 0)

	logs := make(map[string]map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		lineJson := make(map[string]string)
		err := json.Unmarshal(line, &lineJson)
		require.NoError(t, err)
		logs[lineJson["msg"]] = lineJson
	}

	for i := 0; i < count; i++ {
		line := logs[fmt.Sprintf("Run %d", i)]
		key := fmt.Sprintf("key%d", i)
		expectedValue := fmt.Sprintf("value%d", i)
		require.Equal(t, expectedValue, line[key])
	}
}
