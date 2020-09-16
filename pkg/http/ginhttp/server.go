package ginhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/pkg/http/bootstrap"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	server      *http.Server
	beforeFuncs []bootstrap.BeforeServerStartFunc
	afterFuncs  []bootstrap.AfterServerStopFunc
	opts        ServerOptions
	exit        chan os.Signal
}

func NewServer() *Server {
	opts := NewServerOptions(conf.AppConfig)
	s := new(Server)
	s.opts.IdleTimeout = opts.IdleTimeout
	s.opts.ReadTimeout = opts.ReadTimeout
	s.opts.WriteTimeout = opts.WriteTimeout
	s.opts.Addr = opts.Addr
	s.opts.Mode = opts.Mode

	// 设置开发模式
	SetRuntimeMode(opts.Mode)

	handler := gin.New()
	server := &http.Server{
		Addr:         opts.Addr,
		Handler:      handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}
	s.exit = make(chan os.Signal, 2)
	s.server = server
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

// RunHTTPServer provide run http or https protocol.
func (s *Server) RunHTTPServer() error {
	if conf.AppConfig.AutoTLS.Enabled {
		return s.autoTLSServer()
	} else if len(conf.AppConfig.TLS.CertPath) > 0 && len(conf.AppConfig.TLS.KeyPath) > 0 {
		return s.defaultTLSServer()
	} else {
		return s.defaultServer()
	}
}

// Serve serve http request
func (s *Server) defaultServer() error {
	var err error
	for _, fn := range s.beforeFuncs {
		err = fn()
		if err != nil {
			return err
		}
	}

	go s.waitShutdown()
	err = s.server.ListenAndServe()

	for _, fn := range s.afterFuncs {
		fn()
	}
	return err
}

func (s *Server) autoTLSServer() error {
	var g errgroup.Group

	dir := filepath.Join(os.Getenv("HOME"), ".cache", "go-autocert")
	_ = os.MkdirAll(dir, 0700)

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(conf.Config.App.AutoTLS.Host),
		// Cache:      autocert.DirCache(app.config.Core.AutoTLS.Folder),
		Cache: autocert.DirCache(dir),
	}

	g.Go(func() error {
		return http.ListenAndServe(":http", manager.HTTPHandler(http.HandlerFunc(s.redirect)))
	})

	g.Go(func() error {
		var err error
		for _, fn := range s.beforeFuncs {
			err = fn()
			if err != nil {
				return err
			}
		}

		s.server.TLSConfig = &tls.Config{
			GetCertificate: manager.GetCertificate,
			NextProtos:     []string{"http/1.1"}, // disable h2 because Safari :(
		}

		go s.waitShutdown()
		log.Printf("Start to listening the incoming requests on https address")
		err = s.server.ListenAndServeTLS("", "")

		for _, fn := range s.afterFuncs {
			fn()
		}
		return err
	})

	return g.Wait()
}

func (s *Server) defaultTLSServer() error {
	var g errgroup.Group
	g.Go(func() error {
		return http.ListenAndServe(":http", http.HandlerFunc(s.redirect))
	})
	g.Go(func() error {
		var err error
		for _, fn := range s.beforeFuncs {
			err = fn()
			if err != nil {
				return err
			}
		}

		s.server.Addr = fmt.Sprintf("%s:%d", "0.0.0.0", conf.AppConfig.TLS.Port)
		s.server.TLSConfig = &tls.Config{
			NextProtos: []string{"http/1.1"}, // disable h2 because Safari :(
		}

		go s.waitShutdown()
		log.Printf("Start to listening the incoming requests on https address: %d", conf.AppConfig.TLS.Port)
		err = s.server.ListenAndServeTLS(
			conf.Config.App.TLS.CertPath,
			conf.Config.App.TLS.KeyPath,
		)

		for _, fn := range s.afterFuncs {
			fn()
		}
		return err
	})
	return g.Wait()
}

// redirect ...
func (s *Server) redirect(w http.ResponseWriter, req *http.Request) {
	var serverHost = conf.Config.App.Host
	serverHost = strings.TrimPrefix(serverHost, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	req.URL.Scheme = "https"
	req.URL.Host = serverHost

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
}

// Shutdown close http server
func (s *Server) waitShutdown() {
	<-s.exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("GIN-Server", "shutdown http server ...")

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("GIN-Server", "shutdown http server error:%s", err)
	}

	os.Exit(0)

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

func (s *Server) GetServer() *http.Server {
	return s.server
}

func (s *Server) GetGinEngine() *gin.Engine {
	return s.server.Handler.(*gin.Engine)
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
