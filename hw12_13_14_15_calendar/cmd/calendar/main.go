package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service/event"
	storage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/pflag"
)

func main() {
	var configFile string
	pflag.StringVar(&configFile, "cfg", "./configs/config.toml", "Path to configuration file")
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}
	cfg := config.NewConfig(configFile)
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Type)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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
	calendar := event.New(logg, store)

	server := internalhttp.NewServer(cfg, logg, calendar)

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
		if !errors.Is(err, http.ErrServerClosed) {
			logg.Error("http server: " + err.Error())
		}
		cancel()
	}
}
