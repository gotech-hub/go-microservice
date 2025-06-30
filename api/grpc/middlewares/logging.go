package middlewares

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger "go-source/pkg/log"
)

// LoggingInterceptor logs gRPC unary requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	log := logger.GetLogger()
	log.Info().
		Str("method", info.FullMethod).
		Str("request_id", getRequestID(ctx)).
		Msg("gRPC request started")

	// Call the handler
	resp, err := handler(ctx, req)

	// Log response
	duration := time.Since(start)
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		} else {
			statusCode = codes.Unknown
		}
	}

	log.Info().
		Str("method", info.FullMethod).
		Str("request_id", getRequestID(ctx)).
		Dur("duration", duration).
		Str("status", statusCode.String()).
		Msg("gRPC request completed")

	return resp, err
}

// StreamLoggingInterceptor logs gRPC stream requests
func StreamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	log := logger.GetLogger()
	log.Info().
		Str("method", info.FullMethod).
		Str("request_id", getRequestID(ss.Context())).
		Msg("gRPC stream started")

	// Call the handler
	err := handler(srv, ss)

	// Log response
	duration := time.Since(start)
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		} else {
			statusCode = codes.Unknown
		}
	}

	log.Info().
		Str("method", info.FullMethod).
		Str("request_id", getRequestID(ss.Context())).
		Dur("duration", duration).
		Str("status", statusCode.String()).
		Msg("gRPC stream completed")

	return err
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	// TODO: Implement request ID extraction from context
	// For now, return a placeholder
	return "grpc-request"
}
