package models

import (
	"time"

	pb "go-source/api/grpc/proto/gen/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// User model cho internal use
type User struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Phone     string     `json:"phone" db:"phone"`
	Status    UserStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// UserStatus enum
type UserStatus int32

const (
	UserStatusUnspecified UserStatus = iota
	UserStatusActive
	UserStatusInactive
	UserStatusSuspended
)

// String returns string representation of UserStatus
func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusSuspended:
		return "suspended"
	default:
		return "unspecified"
	}
}

// ToProto converts User to protobuf User
func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Status:    pb.UserStatus(u.Status),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

// FromProto converts protobuf User to User
func (u *User) FromProto(protoUser *pb.User) {
	u.ID = protoUser.Id
	u.Name = protoUser.Name
	u.Email = protoUser.Email
	u.Phone = protoUser.Phone
	u.Status = UserStatus(protoUser.Status)
	if protoUser.CreatedAt != nil {
		u.CreatedAt = protoUser.CreatedAt.AsTime()
	}
	if protoUser.UpdatedAt != nil {
		u.UpdatedAt = protoUser.UpdatedAt.AsTime()
	}
}

// CreateUserRequest model
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

// UpdateUserRequest model
type UpdateUserRequest struct {
	ID     string     `json:"id" validate:"required"`
	Name   string     `json:"name"`
	Email  string     `json:"email" validate:"omitempty,email"`
	Phone  string     `json:"phone"`
	Status UserStatus `json:"status"`
}

// ListUsersRequest model
type ListUsersRequest struct {
	Page   int32      `json:"page"`
	Limit  int32      `json:"limit"`
	Search string     `json:"search"`
	Status UserStatus `json:"status"`
}

// ListUsersResponse model
type ListUsersResponse struct {
	Users []*User `json:"users"`
	Total int32   `json:"total"`
	Page  int32   `json:"page"`
	Limit int32   `json:"limit"`
}
