package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"sc-auth-service/config/loader"
	"sc-auth-service/o11y/clog"
	"sc-auth-service/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Application panicked: %v\n", r)
			os.Exit(1)
		}
	}()

	loader.LoadEnv()

	clog.LoadZap()
	defer clog.CloseZap()

	cfg, err := loader.LoadConfig()
	if err != nil {
		clog.Log().Fatal("Failed to load configuration", zap.Error(err))
	}

	serverMgr := server.NewServerManager(cfg)

	err = serverMgr.LoadPostgresDBs()
	if err != nil {
		clog.Log().Fatal("Failed to initialize Postgres databases", zap.Error(err))
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		serverMgr.StartServers()
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	serverMgr.ShutdownServers()

	wg.Wait()
	clog.Log().Info("Exiting application gracefully")
}
