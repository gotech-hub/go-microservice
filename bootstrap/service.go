package bootstrap

import (
	"go-source/internal/services"
	"go-source/pkg/database/redis"
)

type Services struct {
	UserService services.UserService
}

func NewServices(repo *Repositories, redis *redis.Client, clients *Clients) *Services {
	service := &Services{
		UserService: services.NewUserService(),
	}
	return service
}
