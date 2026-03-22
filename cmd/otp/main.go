package main

import (
	"log"
	"runtime"

	"desctop-otp/internal/app"
	"desctop-otp/internal/config"
)

func main() {
	runtime.LockOSThread()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
