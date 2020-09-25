package fiberhttp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/go-impatient/gaia/app/conf"
)

// Option can be applied in server
type OptionFunc func(s *Server)

type Options struct {
	Mode  string `json:"mode"`  // run mode 可选 dev/prod/test
	Addr  string `json:"addr"`  // TCP address to listen on, ":3000" if empty
	Grace bool   `json:"grace"` // grace mode 可选 pprof 为空不使用
}

// DefaultOptions default config
func DefaultOptions() Options {
	return Options{
		Addr: ":4000",
	}
}

// NewOptions set config
func NewOptions(c *conf.App) *Options {
	return &Options{
		Mode:  c.Mode,
		Addr:  fmt.Sprintf("%s:%d", c.Host, c.Port),
		Grace: c.Grace,
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

// App option sets custom Fiber App to Server
func App(app *fiber.App) OptionFunc {
	return func(s *Server) {
		s.app = app
	}
}
