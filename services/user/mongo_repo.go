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

func (u *UserRepo) createDoctor(ctx context.Context, userInput DoctorInput) (*PublicDoctor, error) {
	user := Doctor{
		Name:       userInput.Name,
		ProfilePic: userInput.ProfilePic,
		Email:      userInput.Email,
		Password:   userInput.Password,
		Specialty:  userInput.Specialty,
		Documents:  userInput.Documents,
		Role:       userInput.Role,
	}

	userCollection := u.db.Collection("users")

	result, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return &PublicDoctor{
		ID:         user.ID,
		Name:       user.Name,
		Specialty:  user.Specialty,
		Email:      user.Email,
		Documents:  user.Documents,
		Role:       user.Role,
		ProfilePic: user.ProfilePic,
	}, nil
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

func (u *UserRepo) getDoctorByID(ctx context.Context, doctorId primitive.ObjectID) (*PublicDoctor, error) {
	var doctor PublicDoctor

	userCollection := u.db.Collection("users")

	err := userCollection.FindOne(ctx, bson.M{"_id": doctorId}).Decode(&doctor)

	if err != nil {
		return nil, err
	}

	return &doctor, nil
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
