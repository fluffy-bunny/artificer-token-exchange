package main

import (
	"fmt"

	"echo-starter/internal/runtime"

	"github.com/rs/zerolog/log"
)

var version = "Development"

func main() {
	fmt.Println("Version:" + version)
	DumpPath("./")
	r := runtime.New(NewStartup())
	err := r.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run the application")
	}
}
