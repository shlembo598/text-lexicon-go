package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/pkg/logger/sl"
)

const (
	ctxTimeout = 5
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *sqlx.DB
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{
		echo: echo.New(),
		cfg:  cfg,
		db:   db,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	server := &http.Server{
		Addr:         s.cfg.Server.Port,
		ReadTimeout:  s.cfg.Server.Timeout,
		WriteTimeout: s.cfg.Server.Timeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("server is listening on PORT", slog.String("port", s.cfg.Server.Port))
		if err := s.echo.StartServer(server); err != nil {
			sl.Fatalf("failed to start server: %v", err)
		}
	}()

	go func() {
		slog.Info("starting Debug server PORT", slog.String("port", s.cfg.Server.Port))
		err := http.ListenAndServe(s.cfg.Server.PProfPort, http.DefaultServeMux)
		if err != nil {
			slog.Error("fail to PROF ListenAndServe:", sl.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Server is shutting down...")

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	slog.Info("server Exited Properly")
	return s.echo.Shutdown(ctx)
}
