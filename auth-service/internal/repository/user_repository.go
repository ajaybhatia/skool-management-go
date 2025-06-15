package repository

import (
	"context"
	"time"

	"skool-management/auth-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateRefreshToken(userID primitive.ObjectID, refreshToken string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *UserRepository) GetByRefreshToken(refreshToken string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"refresh_token": refreshToken}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
