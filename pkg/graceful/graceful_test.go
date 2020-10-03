package graceful

import (
	"io/ioutil"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"

)

// Represents a net/http.Handler.
type testServer struct {
	timeout time.Duration
}

// Implementing  net/http.Handler interface.
func (s *testServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(s.timeout)
	rw.Write([]byte("hi"))
}

// Enforcing interface implementation.
func TestInterface(t *testing.T) {
	var _ Graceful = Graceful{}
}

// Test gracefully starting a server.
func TestGraceStart(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:64951",
		Handler: &testServer{timeout: time.Second * 0},
	}

	g := New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	res, err := http.Get("http://localhost:64951/")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	all, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "hi", string(all))

	require.NoError(t, g.Stop())
}

// Test gracefully shuting down a server.
func TestGraceStop(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:64952",
		Handler: &testServer{timeout: time.Second * 10},
	}

	g := New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	require.NoError(t, g.Stop())

	_, err := http.Get("http://localhost:64952/")
	require.Error(t, err)
}

// Test with SIGINT a server.
func TestKill(t *testing.T) {
	server := &http.Server{
		Addr:    "localhost:64953",
		Handler: &testServer{timeout: time.Second * 10},
	}

	g := New(server)
	go func() {
		require.NoError(t, g.Start())
	}()
	time.Sleep(time.Second)

	_, err := http.Get("http://localhost:64953/")
	require.Error(t, err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

