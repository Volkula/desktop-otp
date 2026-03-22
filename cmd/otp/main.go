package main

// Regenerate Windows file icon for Explorer: go generate (requires github.com/akavel/rsrc in PATH).
//
//go:generate go run ../../cmd/genico -out ../../build/windows/app.ico -icongo ../../internal/app/icon.go
//go:generate rsrc -ico ../../build/windows/app.ico -o rsrc.syso

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
