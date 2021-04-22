package cli

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mmcloughlin/professor"
)

var (
	pprofServer *http.Server
)

func startPProf() {
	address := venom.GetString(fmt.Sprintf("debug%vaddress", keyDelimiter))
	pprofServer = professor.NewServer(address)

	log.WithField("address", address).Warn("starting pprof server")

	go func() {
		err := pprofServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("unable to start pprof server")
		}
	}()
}

func stopPProf() {
	if pprofServer != nil {
		err := pprofServer.Close()

		if err != nil {
			log.WithError(err).Error("unable to stop pprof server")
		} else {
			log.Warn("stopped pprof server")
		}
	}

	pprofServer = nil
}

func togglePProf() {
	if pprofServer == nil {
		startPProf()
	} else {
		stopPProf()
	}
}

func autostartPProf() {
	if venom.GetBool(fmt.Sprintf("debug%vautostart", keyDelimiter)) {
		startPProf()
	}
}

func initPProf() {
	pflags := rootCmd.PersistentFlags()
	pflags.String(fmt.Sprintf("debug%vaddress", keyDelimiter), "localhost:10001", "Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.")
	pflags.Bool(fmt.Sprintf("debug%vautostart", keyDelimiter), false, "Whether the golang pprof server is automatically started.")
}
