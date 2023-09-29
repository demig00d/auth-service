package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/demig00d/auth-service/config"
	http2 "github.com/demig00d/auth-service/internal/controller/http"
	"github.com/demig00d/auth-service/internal/controller/http/middleware"
	"github.com/demig00d/auth-service/internal/repository"
	"github.com/demig00d/auth-service/internal/service"
	"github.com/demig00d/auth-service/pkg/logger"
	"github.com/demig00d/auth-service/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(cfg config.Config) {
	logger := logger.NewLogger(cfg.Level, os.Stdin)
	logger.SetPrefix("app run ")

	creds := options.Credential{
		Username: cfg.MongoDB.User,
		Password: cfg.MongoDB.Password,
	}

	mongo, err := mongodb.New(context.Background(), creds, cfg.MongoDB.Host, cfg.MongoDB.Port)
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		if err = mongo.Disconnect(); err != nil {
			panic(err)
		}
	}()

	repo := repository.NewRepository(mongo, "test", "users", logger)

	tokenService := service.NewTokenService(
		repo,
		cfg.AccessToken.LifeTime,
		[]byte(cfg.AccessToken.Secret),
		logger,
	)
	userService := service.NewUserService(tokenService, logger)

	server := http2.NewServer(userService, logger)

	router := middleware.LoggingMiddleware(logger)(server)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	server.Logger.Infof("Server Started on %s:%d\n", "0.0.0.0", cfg.HTTP.Port)
	<-interrupt
	server.Logger.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		server.Logger.Fatalf("Server Shutdown Failed:%+v", err)
	}
	server.Logger.Info("Server Exited Properly")

}
