package main

import (
	"log"

	_ "git.lothric.net/examples/go/gogin/api"

	"git.lothric.net/examples/go/gogin/internal/app/cli"
)

//	@title			GoGin
//	@version		0.2.0
//	@description	GoGin service provides the unified gist storage
//	@BasePath		/api
func main() {

	cmd, err := cli.NewCli()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
