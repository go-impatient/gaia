package request

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReuest(t *testing.T) {

	gzipCompress := func(buf []byte) []byte {
		var data bytes.Buffer
		gw := gzip.NewWriter(&data)
		if _, err := gw.Write(buf); err != nil {
			panic(err)
		}
		gw.Close()
		return data.Bytes()
	}

	t.Run("gzip should work", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Encoding", "gzip")

			respJson := []byte(`{"version":"v0.1.0"}`)
			w.Write(gzipCompress(respJson))
		}))
		defer ts.Close()

		version := struct {
			Version string `json:"version"`
		}{}
		resp, err := Get(ts.URL).Result(&version).Do()
		require.Nil(err)

		assert.True(resp.OK())
		assert.Equal("v0.1.0", version.Version)
	})
}
