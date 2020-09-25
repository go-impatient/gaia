package fiberhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/pkg/errors"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/pkg/http/bootstrap"
)

type Server struct {
	app         *fiber.App
	db          *conf.DB
	opts        Options
	beforeFuncs []bootstrap.BeforeServerStartFunc
	afterFuncs  []bootstrap.AfterServerStopFunc
	exit        chan os.Signal
}

// NewServerWithOptions with options
func NewServerWithOptions(opts Options) *Server {
	s := &Server{
		opts: opts,
	}
	return s
}

// NewServer get server instance
func NewServer(options ...OptionFunc) *Server {
	opts := NewOptions(conf.AppConfig)
	s := &Server{
		app:  fiber.New(),
		exit: make(chan os.Signal, 2),
	}
	s.opts.Mode = opts.Mode
	s.opts.Addr = opts.Addr
	s.opts.Grace = opts.Grace

	for _, o := range options {
		o(s)
	}

	return s
}

func Default() *Server {
	server := fiber.New()
	server.Use(
		requestid.New(),
		logger.New(),
		recover.New(),
		pprof.New(),
	)

	return &Server{
		app: server,
	}
}

// Run runs a web server
func (s *Server) Run(addr string) error {
	return s.app.Listen(addr)
}

// Run runs a tls web server
func (s *Server) RunTls(certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", s.opts.Addr, config)
	if err != nil {
		return nil
	}

	return s.app.Listener(ln)
}

func (s *Server) Shutdown() error {
	if s.app == nil {
		return fmt.Errorf("shutdown: fiber app is not found")
	}
	return s.app.Shutdown()
}

func (s *Server) Router() fiber.Router {
	return s.app
}

// Serve serve http request
func (s *Server) Serve() error {
	var err error
	for _, fn := range s.beforeFuncs {
		err = fn()
		if err != nil {
			return err
		}
	}

	if s.opts.Grace {
		s.app.Use(pprof.New())
	} else {
		signal.Notify(s.exit, os.Interrupt, syscall.SIGTERM)
		go s.waitShutdown()
		log.Printf("Fiber-Server http server start and serve:%v", s.opts.Addr)

		if len(conf.AppConfig.TLS.CertPath) > 0 && len(conf.AppConfig.TLS.KeyPath) > 0 {
			log.Printf("Fiber-Server 'RunTls()'")
			err = s.RunTls(conf.AppConfig.TLS.CertPath, conf.AppConfig.TLS.KeyPath)
		} else {
			log.Printf("Fiber-Server 'ListenAndServe()'")
			err = s.Run(s.opts.Addr)
		}
	}

	for _, fn := range s.afterFuncs {
		fn()
	}
	return err
}

// Shutdown 平滑关闭服务
func (s *Server) waitShutdown() {
	<-s.exit

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Fiber-Server shutdown http server ...")

	err := s.Shutdown()
	if err != nil {
		log.Fatalf("Fiber-Server shutdown http server error:%s", err)
	}

	return
}

// PingServer 服务心跳检查
func (s *Server) PingServer() (err error) {
	maxPingCount := conf.Config.App.MaxPingCount
	for i := 0; i < maxPingCount; i++ {
		// Ping the app by sending a GET request to `/health`.
		url := fmt.Sprintf("%s:%d/sd/health", conf.Config.App.Host, conf.Config.App.Port)
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Fatalf("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	err = errors.New("Cannot connect to the router.")
	return
}

// ConfigureOptions 更新配置
func (s *Server) ConfigureOptions(options ...OptionFunc) {
	for _, o := range options {
		o(s)
	}
}

func (s *Server) AddBeforeServerStartFunc(fns ...bootstrap.BeforeServerStartFunc) {
	for _, fn := range fns {
		s.beforeFuncs = append(s.beforeFuncs, fn)
	}
}

func (s *Server) AddAfterServerStopFunc(fns ...bootstrap.AfterServerStopFunc) {
	for _, fn := range fns {
		s.afterFuncs = append(s.afterFuncs, fn)
	}
}
