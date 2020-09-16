package ginhttp

import (
	"fmt"
	"github.com/go-impatient/gaia/app/conf"
	"time"
)

// ServerOptions http server options
type ServerOptions struct {
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

func DefaultOptions() *ServerOptions {
	return &ServerOptions{
		Mode:         "dev",
		Addr:         ":9090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}

func NewServerOptions(c *conf.App) *ServerOptions {
	return &ServerOptions{
		Mode:         c.Mode,
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}
}
