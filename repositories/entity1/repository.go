package entity

import (
	"context"
	"errors"
	"go-source/pkg/database/mongodb"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type EntityRepository struct {
	*mongodb.Repository[Entity]
}

type IEntityRepository interface {
	Create(ctx context.Context, data *Entity) error
	Get(ctx context.Context, id string) (*Entity, error)
}

var (
	instance *EntityRepository
	once     sync.Once
)

func NewEntityRepository(dbStorage *mongodb.DatabaseStorage) IEntityRepository {
	once.Do(func() {
		instance = &EntityRepository{
			Repository: mongodb.NewRepository[Entity](dbStorage, nil),
		}
	})
	return instance
}

func (r *EntityRepository) Create(ctx context.Context, data *Entity) error {
	_, err := r.CreateOneDocument(ctx, data)
	if !errors.Is(err, nil) {
		return err
	}
	return nil
}

func (r *EntityRepository) Get(ctx context.Context, id string) (*Entity, error) {
	r.ApplyFilters(
		byId(id),
	)
	rs, err := r.FindOneDoc(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
