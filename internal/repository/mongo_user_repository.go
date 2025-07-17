package repository

import (
	"context"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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

func (m *MongoUserRepo) CreateUserTx(tx *gorm.DB, user *models.User) error {
	return nil
}

type Filter struct {
	email     string
	tenant_id string
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

func (m *MongoUserRepo) FindUserByEmailAndTenant(email, tenant_id string) (*models.User, error) {
	filter := &Filter{email: email, tenant_id: tenant_id}
	result := m.collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user *models.User

	result.Decode(&user)

	return user, nil
}

func (m *MongoUserRepo) RemoveUserById(id string) error {
	return nil
}
func (m *MongoUserRepo) RemoveUserByEmail(email string, tenant_id string) error {
	return nil
}
