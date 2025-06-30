package handlers

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-source/api/grpc/models"
	pb "go-source/api/grpc/proto/gen/user"
	"go-source/internal/services"
)

// UserServiceHandler implements pb.UserServiceServer
type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
	userService services.UserService
}

// NewUserServiceHandler creates a new UserServiceHandler
func NewUserServiceHandler(userService services.UserService) *UserServiceHandler {
	return &UserServiceHandler{
		userService: userService,
	}
}

// CreateUser implements pb.UserServiceServer.CreateUser
func (h *UserServiceHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate request
	if req.Name == "" || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name and email are required")
	}

	// Convert to internal model
	createReq := &models.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	// Call service
	user, err := h.userService.CreateUser(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Convert to proto response
	protoUser := user.ToProto()
	return &pb.CreateUserResponse{
		User:    protoUser,
		Message: "User created successfully",
	}, nil
}

// GetUser implements pb.UserServiceServer.GetUser
func (h *UserServiceHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	// Call service
	user, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// Convert to proto response
	protoUser := user.ToProto()
	return &pb.GetUserResponse{
		User: protoUser,
	}, nil
}

// UpdateUser implements pb.UserServiceServer.UpdateUser
func (h *UserServiceHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	// Convert to internal model
	updateReq := &models.UpdateUserRequest{
		ID:     req.Id,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: models.UserStatus(req.Status),
	}

	// Call service
	user, err := h.userService.UpdateUser(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Convert to proto response
	protoUser := user.ToProto()
	return &pb.UpdateUserResponse{
		User:    protoUser,
		Message: "User updated successfully",
	}, nil
}

// DeleteUser implements pb.UserServiceServer.DeleteUser
func (h *UserServiceHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id is required")
	}

	// Call service
	err := h.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

// ListUsers implements pb.UserServiceServer.ListUsers
func (h *UserServiceHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Convert to internal model
	listReq := &models.ListUsersRequest{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
		Status: models.UserStatus(req.Status),
	}

	// Call service
	result, err := h.userService.ListUsers(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	// Convert to proto response
	var protoUsers []*pb.User
	for _, user := range result.Users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}, nil
}

// HealthCheck implements pb.UserServiceServer.HealthCheck
func (h *UserServiceHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status:    "OK",
		Message:   "User service is healthy",
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}
