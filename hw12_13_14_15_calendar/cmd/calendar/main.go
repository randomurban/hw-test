package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	implementation "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/api/event"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/server/grpcserver"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/server/httpserver"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service/event"
	storage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/pflag"
)

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}
	cfg := config.NewConfig(configFile)
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Type)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var store storage.EventStorage
	switch cfg.Store.StoreType {
	case config.StoreTypeSQL:
		store = sqlstorage.New()
		err := store.Connect(ctx, cfg.DB.DSN)
		if err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			os.Exit(1)
		}
	case config.StoreTypeMemory:
		store = memorystorage.New()
	default:
		logg.Error("Unknown store type: " + cfg.Store.StoreType)
	}
	defer cancel()
	service := event.New(logg, store)

	impl := implementation.NewImplementation(service)

	httpSrv := httpserver.NewServer(cfg, logg, impl)
	grpcSrv := grpcserver.NewServer(cfg, logg, impl)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		logg.Info("http calendar is running...")

		if err := httpSrv.Start(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logg.Error("http httpSrv: " + err.Error())
			}
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		logg.Info("grpc calendar is running...")

		if err := grpcSrv.Start(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logg.Error("http httpSrv: " + err.Error())
			}
			cancel()
		}
	}()
	wg.Wait()
}
