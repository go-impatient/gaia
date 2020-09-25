package fiberhttp

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Server_New(t *testing.T) {
	app := fiber.New()
	s := NewServer(App(app))

	assert.Equal(t, app, s.app)
	assert.Equal(t, ":4000", s.opts.Addr)
}

func Test_Server_Default(t *testing.T) {
	s := Default()

	require.NotNil(t, s.app)
}

func Test_Server_Run(t *testing.T) {
	s := NewServer()

	go func() {
		time.Sleep(time.Millisecond * 100)
		assert.NoError(t, s.app.Shutdown())
	}()

	assert.NoError(t, s.Run(""))
}

func Test_Server_RunTls(t *testing.T) {
	s := NewServer()

	t.Run("invalid addr", func(t *testing.T) {
		assert.NotNil(t, s.RunTls("./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})

	t.Run("invalid ssl info", func(t *testing.T) {
		assert.NotNil(t, s.RunTls("./.github/README.md", "./.github/README.md"))
	})

	t.Run("with ssl", func(t *testing.T) {
		go func() {
			time.Sleep(time.Millisecond * 100)
			assert.NoError(t, s.app.Shutdown())
		}()

		assert.NoError(t, s.RunTls("./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})
}

func Test_Server_Shutdown(t *testing.T) {
	require.NotNil(t, (&Server{}).Shutdown())
	require.Nil(t, NewServer().Shutdown())
}

func Test_Server_Router(t *testing.T) {
	require.Nil(t, (&Server{}).Router())
	require.NotNil(t, NewServer().Router())
}
