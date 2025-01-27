package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) createUser(ctx context.Context, userInput UserInput) (*PublicUser, error) {
	user := User{
		Name:       userInput.Name,
		ProfilePic: userInput.ProfilePic,
		Email:      userInput.Email,
		Password:   userInput.Password,
		Role:       userInput.Role,
	}

	userCollection := u.db.Collection("users")

	result, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user.ToPublic(), nil
}

func (u *UserRepo) getUserByID(ctx context.Context, userId primitive.ObjectID) (*PublicUser, error) {
	var user User
	userCollection := u.db.Collection("users")

	err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user.ToPublic(), nil
}

func (u *UserRepo) getUserByEmail(ctx context.Context, userEmail string) (*User, error) {
	var user User
	userCollection := u.db.Collection("users")

	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
