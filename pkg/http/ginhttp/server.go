package ginhttp

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

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/pkg/http/bootstrap"
	"github.com/go-impatient/gaia/pkg/http/ginhttp/grace"
)

type Server struct {
	app         *http.Server
	beforeFuncs []bootstrap.BeforeServerStartFunc
	afterFuncs  []bootstrap.AfterServerStopFunc
	opts        Options
	exit        chan os.Signal
}

// NewServerWithOptions with options
func NewServerWithOptions(opts Options) *Server {
	s := &Server{
		opts: opts,
	}
	return s
}

func NewServer(options ...OptionFunc) *Server {
	opts := NewServerOptions(conf.AppConfig)
	s := new(Server)
	s.opts.IdleTimeout = opts.IdleTimeout
	s.opts.ReadTimeout = opts.ReadTimeout
	s.opts.WriteTimeout = opts.WriteTimeout
	s.opts.Addr = opts.Addr
	s.opts.Mode = opts.Mode
	handler := gin.New()
	server := &http.Server{
		Addr:         opts.Addr,
		Handler:      handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}

	s.app = server
	s.exit = make(chan os.Signal, 2)

	for _, o := range options {
		o(s)
	}

	return s
}

// SetRuntimeMode 设置开发模式
func SetRuntimeMode(mode string) {
	switch mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("unknown mode")
	}
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
		grace.Start(s.opts.Addr, s.app)
	} else {
		signal.Notify(s.exit, os.Interrupt, syscall.SIGTERM)
		go s.waitShutdown()
		log.Printf("Gin-Server http server start and serve:%v", s.opts.Addr)
		if conf.AppConfig.AutoTLS.Enabled {
			log.Printf("Gin-Server 'RunAutoTls()'")
			err = s.RunAutoTls(conf.AppConfig.AutoTLS.Folder, conf.AppConfig.AutoTLS.Host)
		} else if len(conf.AppConfig.TLS.CertPath) > 0 && len(conf.AppConfig.TLS.KeyPath) > 0 {
			log.Printf("Gin-Server 'RunTls()'")
			err = s.RunTls(conf.AppConfig.TLS.CertPath, conf.AppConfig.TLS.KeyPath)
			//} else if s.app.TLSConfig == nil {
			//	log.Printf("Gin-Server 'Run()'")
			//	err = s.Run(s.opts.Addr)
		} else {
			log.Printf("Gin-Server 'ListenAndServe()'")
			err = s.Run()
		}
	}

	for _, fn := range s.afterFuncs {
		fn()
	}
	return err
}

// Run runs a web server
func (s *Server) Run() error {
	return s.app.ListenAndServe()
}

// Run  LetsEncrypt HTTPS servers.
func (s *Server) RunDefTls(addr ...string) error {
	r := s.GetGinEngine()
	return http.Serve(autocert.NewListener(addr...), r)
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
		NextProtos:   []string{"http/1.1"}, // disable h2 because Safari :(
	}

	s.app.TLSConfig = config

	return s.app.ListenAndServeTLS("", "")
}

// Run runs a auto tls web server, host = ("example1.com", "example2.com")
func (s *Server) RunAutoTls(folder string, host ...string) error {

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host...),
		Cache:      autocert.DirCache(folder),
	}

	s.app.TLSConfig = m.TLSConfig()
	s.app.Addr = ":https"

	go http.ListenAndServe(":http", m.HTTPHandler(http.HandlerFunc(s.redirect)))
	return s.app.ListenAndServeTLS("", "")
}

func (s *Server) Router() *gin.Engine {
	return s.GetGinEngine()
}

func (s *Server) GetGinEngine() *gin.Engine {
	return s.app.Handler.(*gin.Engine)
}

func (s *Server) GetServer() *http.Server {
	return s.app
}

// Shutdown 平滑关闭服务
func (s *Server) waitShutdown() {
	<-s.exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("GIN-Server shutdown http server ...")

	err := s.app.Shutdown(ctx)

	if err != nil {
		log.Fatalf("GIN-Server shutdown http server error:%s", err)
	}

	return
}

func (s *Server) redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
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
