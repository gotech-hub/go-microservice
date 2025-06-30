package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-source/api/grpc/middlewares"
	"go-source/internal/services"
	logger "go-source/pkg/log"
)

// Server represents the gRPC server
type Server struct {
	server      *grpc.Server
	port        int
	userService *services.UserService
}

// NewServer creates a new gRPC server
func NewServer(port int, userService *services.UserService) *Server {
	// Create gRPC server with middleware
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares.LoggingInterceptor),
		grpc.StreamInterceptor(middlewares.StreamLoggingInterceptor),
	)

	return &Server{
		server:      server,
		port:        port,
		userService: userService,
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	// Register services
	s.registerServices()

	// Enable reflection for grpcurl
	reflection.Register(s.server)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logger.GetLogger().Info().Msgf("gRPC server starting on port %d", s.port)

	// Start server
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	logger.GetLogger().Info().Msg("Stopping gRPC server...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.GetLogger().Warn().Msg("gRPC server force shutdown")
		s.server.Stop()
	case <-done:
		logger.GetLogger().Info().Msg("gRPC server stopped gracefully")
	}
}

// registerServices registers all gRPC services
func (s *Server) registerServices() {
	// Register user service
	// userHandler := handlers.NewUserServiceHandler(s.userService)
	// Note: We'll need to import the generated pb package after proto generation
	// pb.RegisterUserServiceServer(s.server, userHandler)
}
