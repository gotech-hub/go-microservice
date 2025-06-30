package services

import (
	"context"
	"time"

	"go-source/api/grpc/models"
	"go-source/pkg/utils"
)

// UserService interface
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, req *models.ListUsersRequest) (*models.ListUsersResponse, error)
}

// userService implementation
type userService struct {
	// Add repository dependency here when needed
}

// NewUserService creates a new UserService
func NewUserService() UserService {
	return &userService{}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Generate unique ID
	id := utils.GenerateID()

	// Create user
	user := &models.User{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    models.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// TODO: Save to database using repository
	// For now, return mock data
	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	// TODO: Get from database using repository
	// For now, return mock data
	if id == "" {
		return nil, utils.ErrNotFound
	}

	return &models.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
		Status:    models.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, req *models.UpdateUserRequest) (*models.User, error) {
	// TODO: Update in database using repository
	// For now, return mock data
	user := &models.User{
		ID:        req.ID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	// TODO: Delete from database using repository
	// For now, just return success
	return nil
}

// ListUsers retrieves a list of users with pagination
func (s *userService) ListUsers(ctx context.Context, req *models.ListUsersRequest) (*models.ListUsersResponse, error) {
	// TODO: Get from database using repository with pagination
	// For now, return mock data
	users := []*models.User{
		{
			ID:        "1",
			Name:      "User 1",
			Email:     "user1@example.com",
			Phone:     "+1234567890",
			Status:    models.UserStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			Name:      "User 2",
			Email:     "user2@example.com",
			Phone:     "+1234567891",
			Status:    models.UserStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return &models.ListUsersResponse{
		Users: users,
		Total: int32(len(users)),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}
