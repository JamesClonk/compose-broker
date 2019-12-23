package main

import (
	"net/http"

	"github.com/JamesClonk/compose-broker/broker"
	"github.com/JamesClonk/compose-broker/config"
	"github.com/JamesClonk/compose-broker/env"
	"github.com/JamesClonk/compose-broker/log"
)

func main() {
	port := env.Get("PORT", "8080")

	log.Infoln("port:", port)
	log.Infoln("log level:", config.Get().LogLevel)
	log.Infoln("broker username:", config.Get().Username)
	log.Infoln("api url:", config.Get().API.URL)
	log.Infoln("api default datacenter:", config.Get().API.DefaultDatacenter)

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, broker.NewRouter(config.Get())))
}
