package runner

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/repository"
	"backend/internal/router"
	"backend/internal/service"
	"backend/pkg"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var logger *slog.Logger

func init() {
	logger = pkg.NewLogger(os.Stdout, slog.LevelDebug, "backend")
}

func Run() {
	logger.Info("starting application")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		return
	}

	connect, err := config.Connect(cfg)
	if err != nil {
		logger.Error("connect database", "error", err)
		return
	}
	defer connect.Close()

	bookingRepository := repository.NewBookingRepository(connect, logger)
	movieRepository := repository.NewMovieRepository(connect, logger)

	bookingService := service.NewBookingService(bookingRepository, logger)
	movieService := service.NewMovieService(movieRepository, logger)

	bookingHandler := api.NewBookingHandler(bookingService, logger)
	movieHandler := api.NewMovieHandler(movieService, logger)

	serverMux := router.Create(bookingHandler, movieHandler, logger)
	startGraceful(cfg, serverMux)
}

func start(conf *config.Configuration, router *mux.Router) error {
	return http.ListenAndServe(conf.Server.GetAddr(), router)
}

func startGraceful(conf *config.Configuration, handler http.Handler) {
	server := &http.Server{
		Addr:    conf.Server.GetAddr(),
		Handler: handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server started", "port", conf.Server.GetAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-stop
	logger.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}
