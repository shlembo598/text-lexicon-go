package main

import (
	"fmt"
	"log/slog"

	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/internal/server"
	"github.com/shlembo598/text-lexicon-go/pkg/db/postgres"
	"github.com/shlembo598/text-lexicon-go/pkg/logger/sl"
)

func main() {
	cfg := config.MustLoad()

	sl.SetupLogger(cfg.Env)

	slog.Info(
		"Config",
		slog.String("AppVersion", cfg.Server.AppVersion),
		slog.String("Env", cfg.Env),
		slog.String("Mode", cfg.Server.Mode),
	)

	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		sl.Fatalf("Postgresql init: %s", err)
	} else {
		slog.Info("Postgres connected", slog.String("status", fmt.Sprint(db.Stats())))
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Error("failed to close storage", sl.Err(err))
		}
	}()

	s := server.NewServer(cfg, db)
	if err = s.Run(); err != nil {
		sl.Fatalf("failed to run server", err)
	}

}
