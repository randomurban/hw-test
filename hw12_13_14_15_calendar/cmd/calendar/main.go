package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service"
	storage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/pflag"
)

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}
	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var store storage.EventStorage
	switch config.Store.StoreType {
	case StoreTypeSQL:
		store = sqlstorage.New()
		err := store.Connect(ctx, config.DB.DSN)
		if err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			os.Exit(1)
		}
	case StoreTypeMemory:
		store = memorystorage.New()
	default:
		logg.Error("Unknown store type: " + config.Store.StoreType)
	}
	calendar := service.New(logg, store)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()
		logg.Info("calendar is stopping...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	// sample1(ctx, store, logg)

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
