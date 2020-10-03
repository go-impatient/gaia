package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/pkg/server/bootstrap"
	"github.com/go-impatient/gaia/pkg/server/grace"
)

type Server struct {
	app         http.Handler
	beforeFuncs []bootstrap.BeforeServerStartFunc
	afterFuncs  []bootstrap.AfterServerStopFunc
	opts        Options
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

	for _, o := range options {
		o(s)
	}

	return s
}

// Serve serve http request
func (s *Server) Serve(ctx context.Context) error {
	var err error
	for _, fn := range s.beforeFuncs {
		err = fn()
		if err != nil {
			return err
		}
	}

	if conf.AppConfig.TLS.Enabled {
		return s.RunTls(ctx)
	} else if conf.AppConfig.AutoTLS.Enabled {
		return s.RunAutoTls(ctx)
	} else {
		return s.Run(ctx)
	}

	for _, fn := range s.afterFuncs {
		err = fn()
		if err != nil {
			return err
		}
	}

	return err
}

// Run runs a web server
func (s *Server) Run(ctx context.Context) error {
	var g errgroup.Group
	opts := NewServerOptions(conf.AppConfig)
	s1 := &http.Server{
		Addr:  opts.Addr,
		Handler: s.app,
		ReadTimeout: s.opts.ReadTimeout,
		WriteTimeout: s.opts.WriteTimeout,
		IdleTimeout: s.opts.IdleTimeout,
	}

	if s.opts.Grace {
		grace.Start(s.opts.Addr, s1)
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("server shutdown http server ...")
			return s1.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		return s1.ListenAndServe()
	})

	return g.Wait()
}

// Run runs a tls web server
func (s *Server) RunTls(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(s.redirect),
		ReadTimeout: s.opts.ReadTimeout,
		WriteTimeout: s.opts.WriteTimeout,
		IdleTimeout: s.opts.IdleTimeout,
	}
	s2 := &http.Server{
		Addr:    ":https",
		Handler: s.app,
		ReadTimeout: s.opts.ReadTimeout,
		WriteTimeout: s.opts.WriteTimeout,
		IdleTimeout: s.opts.IdleTimeout,
	}

	g.Go(func() error {
		return s1.ListenAndServe()
	})

	g.Go(func() error {
		return s2.ListenAndServeTLS(conf.AppConfig.TLS.CertPath, conf.AppConfig.TLS.KeyPath)
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("server shutdown http/https server ...")
			s1.Shutdown(ctx)
			s2.Shutdown(ctx)
			return nil
		}
	})

	return g.Wait()
}

// Run runs a auto tls web server, hosts = ("example1.com", "example2.com")
func (s *Server) RunAutoTls(ctx context.Context) error {
	var g errgroup.Group
	certManager := autocert.Manager{
		Cache: autocert.DirCache(conf.AppConfig.AutoTLS.Folder),
		Prompt: func(tosURL string) bool {
			return conf.AppConfig.AutoTLS.AcceptTos
		},
		HostPolicy: autocert.HostWhitelist(conf.AppConfig.AutoTLS.Hosts...),
	}
	httpHandler := certManager.HTTPHandler(http.HandlerFunc(s.redirect))
	s1 := &http.Server{
		Addr:    ":http",
		Handler: httpHandler,
		ReadTimeout: s.opts.ReadTimeout,
		WriteTimeout: s.opts.WriteTimeout,
		IdleTimeout: s.opts.IdleTimeout,
	}
	config := &tls.Config{
		GetCertificate: certManager.GetCertificate,
		MinVersion:     tls.VersionTLS12,
		NextProtos:     []string{"http/1.1"}, // []string{"h2", "http/1.1"}, disable h2 because Safari :(
	}
	s2 := &http.Server{
		Addr:      ":https",
		Handler:   s.app,
		ReadTimeout: s.opts.ReadTimeout,
		WriteTimeout: s.opts.WriteTimeout,
		IdleTimeout: s.opts.IdleTimeout,
		TLSConfig: config,
	}
	g.Go(func() error {
		return s1.ListenAndServe()
	})
	g.Go(func() error {
		return s2.ListenAndServeTLS("", "")
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("server shutdown http/https server ...")
			s1.Shutdown(ctx)
			s2.Shutdown(ctx)
			return nil
		}
	})
	return g.Wait()
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
