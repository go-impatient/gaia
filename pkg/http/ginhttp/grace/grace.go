package grace

import (
	"net/http"
	"runtime/debug"

	"github.com/go-impatient/gaia/pkg/pprof"
	"github.com/google/martian/log"
	"github.com/jpillora/overseer"
	"github.com/spf13/cast"
)

var server = &http.Server{}
var addresses = make([]string, 0)

func Start(addr string, s *http.Server) {
	addresses = append(addresses, addr)
	server = s
	if pprof.PprofPort != "" {
		addresses = append(addresses, pprof.PprofPort)
	}

	oversee()
}

func oversee() {
	overseer.Run(overseer.Config{
		Program:   prog,
		Addresses: addresses,
		//Fetcher: &fetcher.File{Path: "my_app_next"},
		Debug: true, //display log of overseer actions
	})

}

func prog(state overseer.State) {
	log.Infof("Program", "app (%s) listening...\n", state.ID)
	if len(addresses) > 1 {
		for k, v := range addresses {
			if v == pprof.PprofPort {
				go pprof.Start(state.Listeners[k])
			}
		}

	}
	err := server.Serve(state.Listener)
	if err != nil {
		log.Errorf("ServerError", "Unhandled error: %v\n stack:%v", err.Error(), cast.ToString(debug.Stack()))
	}
	log.Infof("Program", "app (%s) exiting...\n", state.ID)
}
