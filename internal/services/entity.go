package services

import (
	"context"
	"go-source/internal/domains"
	"go-source/pkg/adapters"
	entity "go-source/repositories/entity1"
	"sync"
)

type EntityService struct {
	entityRepository entity.IEntityRepository
}

type IEntityService interface {
	Get(ctx context.Context, Id string) (*domains.Entity, error)
}

var (
	entityInst *EntityService
	onceEntity sync.Once
)

func NewEntityService(entityRepository entity.IEntityRepository) IEntityService {
	onceEntity.Do(func() {
		entityInst = &EntityService{
			entityRepository: entityRepository,
		}
	})
	return entityInst
}

func (s *EntityService) Get(ctx context.Context, Id string) (*domains.Entity, error) {
	// Convert Id to primitive.ObjectID
	rs, err := s.entityRepository.Get(ctx, Id)
	if err != nil {
		return nil, err
	}
	r := adapters.AdapterEntity{}.ConvertServToRepo(rs)
	return r, nil
}
