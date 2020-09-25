package ginhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-impatient/gaia/app/conf"
)

// Option can be applied in server
type OptionFunc func(s *Server)

// Options http server options
type Options struct {
	// run mode 可选 dev/prod/test
	Mode string `json:"mode"`
	// TCP address to listen on, ":http" if empty
	Addr string `json:"addr"`
	//grace mode 可选graceful/oversea 为空不使用
	Grace bool `json:"grace"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `json:"read_timeout"`
	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a create
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration `json:"write_timeout"`
	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration `json:"idle_timeout"`
}

func DefaultOptions() *Options {
	return &Options{
		Mode:         "dev",
		Addr:         ":4000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}

func NewServerOptions(c *conf.App) *Options {
	return &Options{
		Mode:         c.Mode,
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}
}

func Addr(a string) OptionFunc {
	return func(s *Server) {
		s.opts.Addr = a
	}
}

func Mode(a string) OptionFunc {
	return func(s *Server) {
		s.opts.Mode = a
	}
}

func Grace(a bool) OptionFunc {
	return func(s *Server) {
		s.opts.Grace = a
	}
}

func ReadTimeout(a time.Duration) OptionFunc {
	return func(s *Server) {
		s.opts.ReadTimeout = a
	}
}

func WriteTimeout(a time.Duration) OptionFunc {
	return func(s *Server) {
		s.opts.WriteTimeout = a
	}
}

func IdleTimeout(a time.Duration) OptionFunc {
	return func(s *Server) {
		s.opts.IdleTimeout = a
	}
}

// App option sets custom Gin App to Server
func App(app *http.Server) OptionFunc {
	return func(s *Server) {
		s.app = app
	}
}
