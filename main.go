package main

import (
	"authserver/common"
	"authserver/dependencies"
	"authserver/server"
	"log"

	"authserver/config"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(common.ChainError("error initing config", err))
	}

	serverRunner := server.CreateHTTPServerRunner(dependencies.ResolveRouterFactory())
	log.Fatal(serverRunner.Run())
}
