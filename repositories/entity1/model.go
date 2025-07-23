package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Entity struct {
	Id        string    `bson:"_id,omitempty"`
	Status    string    `bson:"status"` // pending, accepted, blocked
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (Entity) CollectionName() string {
	return ColEntity
}

func (Entity) IndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}
