package server

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
	ReadTimeout time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
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
func App(h http.Handler) OptionFunc {
	return func(s *Server) {
		s.app = h
	}
}
