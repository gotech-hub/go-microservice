package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"go-source/api/grpc"
	"go-source/api/http"
	"go-source/api/ws"
	"go-source/bootstrap"
	"go-source/config"
	"go-source/pkg/constant"
	"go-source/pkg/database/redis"
	logger "go-source/pkg/log"
)

func main() {
	// Initialize logging system with service name
	logger.InitLog("go-source")
	log := logger.GetLogger()
	log.Info().Msgf("Start %s services", constant.ServiceName)

	// Load application configuration from environment variables
	config, err := config.LoadConfig()
	if err != nil {
		logger.GetLogger().Fatal().Msgf("Failed to load configuration: %v", err)
		return
	}

	// Set health check status to true for service discovery
	http.SetHealthCheck(true)
	e := echo.New()

	// Setup context with graceful shutdown signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize Redis connection
	redisClient, err := redis.ConnectRedis(ctx, &config.RedisConfig)
	if err != nil {
		log.Fatal().Msgf("Connect redis failed: %s", err)
	}

	// Initialize application dependencies following clean architecture pattern
	storage := bootstrap.NewDatabaseConnection(ctx)
	clients := bootstrap.NewClients()
	repositores := bootstrap.NewRepositories(storage.Connection)
	services := bootstrap.NewServices(repositores, redisClient, clients)
	handlers := bootstrap.NewHandlers(services)

	// Start HTTP server with configured handlers
	srv := http.NewHttpServe(handlers)
	go func() {
		srv.Start(e)
	}()

	// Start WebSocket server
	wsServer := ws.NewWSServer(handlers)
	go func() {
		wsServer.Start(e)
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer(9090, &services.UserService)
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatal().Msgf("gRPC server failed: %v", err)
		}
	}()

	log.Info().Msg("HTTP, WebSocket, and gRPC servers started successfully")

	// Wait for termination signal for graceful shutdown
	<-ctx.Done()
	http.SetHealthCheck(false)
	ws.SetWSHealthCheck(false)

	// Allow 15 seconds for active connections to close
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Stop gRPC server
	grpcServer.Stop()

	// Stop WebSocket server
	if err := wsServer.Stop(shutdownCtx); err != nil {
		log.Error().Msgf("WebSocket server shutdown error: %v", err)
	}

	// Stop HTTP server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msgf("Force shutdown services")
	}
}
