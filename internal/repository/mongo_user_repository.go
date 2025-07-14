package repository

import (
	"context"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(col *mongo.Collection) UserRepository {
	return &MongoUserRepo{collection: col}
}

func (m *MongoUserRepo) CreateUser(user *models.User) error {
	_, err := m.collection.InsertOne(context.TODO(), user)
	return err
}

type Filter struct {
	email string
}

func (m *MongoUserRepo) FindUserByEmail(email string) (*models.User, error) {
	filter := &Filter{email: email}
	result := m.collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user *models.User

	result.Decode(&user)

	return user, nil
}
