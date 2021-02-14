package xonce

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnce_Do(t *testing.T) {
	var once Once
	var conn net.Conn
	var err error

	err = once.Do(func() error {
		conn, err = net.DialTimeout("tcp", "google.com:80", time.Second)
		fmt.Printf("google.com连接失败: %s", err.Error())
		return err
	})

	assert.Nil(t, err)
	assert.NotNil(t, conn)

	if conn != nil {
		conn.Write([]byte("GET / HTTP/1.1\r\nHost: google.com\r\n Accept: */*\r\n\r\n"))
		io.Copy(os.Stdout, conn)
	}
}

func TestOnce_Done(t *testing.T) {
	var flag Once

	assert.False(t, flag.Done())

	flag.Do(func() error {
		time.Sleep(time.Second)
		return nil
	})

	assert.True(t, flag.Done())
}
