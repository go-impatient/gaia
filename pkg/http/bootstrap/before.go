package bootstrap

import (
	"github.com/go-impatient/gaia/pkg/limit"
	"github.com/go-impatient/gaia/pkg/pprof"
)

type BeforeServerStartFunc func() error

func GrowMaxFd() BeforeServerStartFunc {
	return func() error {
		return limit.GrowToMaxFdLimit()
	}
}

func InitPprof() BeforeServerStartFunc {
	return func() error {
		go pprof.Pprof()
		return nil
	}
}
