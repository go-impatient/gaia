package graceful

import (
	"context"
	"net/http"
	"time"

	"github.com/ory/graceful"
)

// Holds *http.Server that needs to be stopped graceful.
type Graceful struct {
	server *http.Server
}

// Initialize Graceful
func New(srv *http.Server) *Graceful {
	g := &Graceful{
		server: graceful.WithDefaults(srv),
	}
	return g
}

// Starts a http.Server that shuts down on SIGINT or SIGTERM.
func (g *Graceful) Start() error {
	return graceful.Graceful(g.start, g.stop)
}

// Callback executed by graceful to start the server.
func (g *Graceful) start() error {
	return g.server.ListenAndServe()
}

// Callback to be called on SIGINT or SIGTERM.
func (g *Graceful) stop(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}

// Shuts down the server.
func (g *Graceful) Stop() error {
	timer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return g.server.Shutdown(timer)
}

