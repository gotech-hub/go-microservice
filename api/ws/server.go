package ws

import (
	"context"
	"errors"
	"fmt"
	"go-source/api/ws/handlers"
	"go-source/bootstrap"
	"go-source/config"
	logger "go-source/pkg/log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	wsHealthCheck bool
	wsMu          sync.RWMutex
)

func SetWSHealthCheck(status bool) {
	wsMu.Lock()
	defer wsMu.Unlock()
	wsHealthCheck = status
}

type WSServInterface interface {
	Start(e *echo.Echo)
	Stop(ctx context.Context) error
}

type WSServer struct {
	handlers *bootstrap.Handlers
	wsEcho   *echo.Echo
}

func NewWSServer(handlers *bootstrap.Handlers) *WSServer {
	return &WSServer{
		handlers: handlers,
	}
}

func (app *WSServer) Start(e *echo.Echo) {
	log := logger.GetLogger()
	wsPort := config.GetInstance().WSPort
	if wsPort == 0 {
		wsPort = 8081 // Default WebSocket port
	}

	// Create a new Echo instance for WebSocket server
	app.wsEcho = echo.New()

	// Setup WebSocket routes
	wsHandler := handlers.NewWSHandler(app.handlers)
	app.wsEcho.GET("/ws", wsHandler.HandleWebSocket)
	app.wsEcho.GET("/ws/chat", wsHandler.HandleChatWebSocket)
	app.wsEcho.GET("/ws/notifications", wsHandler.HandleNotificationWebSocket)

	go func() {
		err := app.wsEcho.Start(fmt.Sprintf(":%d", wsPort))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("can't start WebSocket server: %v", err)
		}
	}()

	log.Info().Msgf("WebSocket server started on port %d", wsPort)
}

func (app *WSServer) Stop(ctx context.Context) error {
	if app.wsEcho != nil {
		return app.wsEcho.Shutdown(ctx)
	}
	return nil
}
