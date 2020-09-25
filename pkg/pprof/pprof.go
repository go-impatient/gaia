package pprof

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Pprof() {
	grace := viper.GetBool("grace")
	if grace == true {
		InitPort()
	} else {
		go pprofstartWithoutGrace()
	}
}

func pprofstartWithoutGrace() {
	enable := viper.Get("enable")
	if enable != "true" {
		return
	}

	port := viper.GetString("port")
	if len(port) <= 0 {
		log.Printf("Pprof", "pprof port:%s format wrong", port)
		return
	}
	log.Printf("Pprof", "open pprof on port:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var PprofPort string

var Pserver = &http.Server{Addr: PprofPort}

func InitPort() {
	enable := viper.Get("enable")
	if enable != "true" {
		return
	}

	PprofPort = ":" + viper.GetString("port")
	if len(PprofPort) <= 0 {
		log.Printf("Pprof", "pprof port:%s format wrong", PprofPort)
		return
	}

}

func Start(l net.Listener) {
	err := Pserver.Serve(l)
	if err != nil {
		log.Printf("ServerError", "Unhandled error: %v\n stack:%v", err.Error(), cast.ToString(debug.Stack()))
	}
}
