package main

import (
	"log"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/dependencies"
	"github.com/mhogar/amber/server"

	"github.com/mhogar/amber/config"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(common.ChainError("error initing config", err))
	}

	serverRunner := server.CreateHTTPServerRunner(dependencies.ResolveRouterFactory())
	log.Fatal(serverRunner.Run())
}
