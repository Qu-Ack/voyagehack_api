package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo *UserRepo
}

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{
		repo: userRepo,
	}
}

func (u *UserService) CreateRoot(ctx context.Context, userInput UserInput, requester PublicUser) (*PublicUser, error) {

	_, err := u.repo.getUserByEmail(ctx, userInput.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	if requester.Role != RoleRoot {
		return nil, errors.New("Only root users can create other root users")
	}
	userInput.Role = RoleRoot

	// password hashing feature not working properly
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	//if err != nil {
	//return nil, err
	//}
	//userInput.Password = string(hashedPassword)

	return u.repo.createUser(ctx, userInput)
}

func (u *UserService) CreateDoctor(ctx context.Context, doctorInput DoctorInput, requester PublicUser) (*PublicDoctor, error) {
	_, err := u.repo.getUserByEmail(ctx, doctorInput.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}
	if requester.Role != RoleRoot {
		return nil, errors.New("Only root users can create other root users")
	}

	doctorInput.Role = RoleDoctor

	return u.repo.createDoctor(ctx, doctorInput)
}

func (u *UserService) CreateStaff(ctx context.Context, userInput UserInput, requester PublicUser) (*PublicUser, error) {
	_, err := u.repo.getUserByEmail(ctx, userInput.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}
	if requester.Role != RoleRoot {
		return nil, errors.New("Only Root users can create other root users")
	}

	userInput.Role = RoleStaff

	// password hashing feature not working properly
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	//if err != nil {
	//return nil, err
	//}

	//userInput.Password = string(hashedPassword)

	return u.repo.createUser(ctx, userInput)
}

func (u *UserService) Login(ctx context.Context, loginInput LoginInput) (*AuthResponse, error) {
	user, err := u.repo.getUserByEmail(ctx, loginInput.Email)

	if err != nil {
		return nil, errors.New("user doesn't exist")
	}

	if loginInput.Password != user.Password {
		return nil, errors.New("passwords don't match")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour).Unix(),
		"id":    user.ID,
		"role":  user.Role,
		"email": user.Email,
	})
	s, err := t.SignedString([]byte("tryandbruteforcethisbitch"))

	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		Token: s,
		User:  *user.ToPublic(),
	}, nil

}

func (u *UserService) Me(ctx context.Context, userId string) (*PublicUser, error) {

	ObjectUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.getUserByID(ctx, ObjectUserId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) MeDoctor(ctx context.Context, doctorId string) (*PublicDoctor, error) {
	DoctorObjectId, err := primitive.ObjectIDFromHex(doctorId)

	if err != nil {
		return nil, err
	}

	doctor, err := u.repo.getDoctorByID(ctx, DoctorObjectId)

	if err != nil {
		return nil, err
	}

	return doctor, nil
}

func (u *UserService) MeByEmail(ctx context.Context, userEmail string) (*PublicUser, error) {
	user, err := u.repo.getUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return &PublicUser{
		ID:         user.ID.Hex(),
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		ProfilePic: user.ProfilePic,
	}, nil

}

func (u *UserService) CreateTestUser(ctx context.Context, userInput UserInput) (*PublicUser, error) {

	_, err := u.repo.getUserByEmail(ctx, userInput.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	// hashing feature not working for now
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 14)
	//fmt.Println(string(hashedPassword))
	//if err != nil {
	//return nil, err
	//}
	//userInput.Password = string(hashedPassword)

	return u.repo.createUser(ctx, userInput)
}

func (u *UserService) CreatePatient(ctx context.Context, userInput UserInput) (*PublicUser, error) {
	_, err := u.repo.getUserByEmail(ctx, userInput.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	return u.repo.createUser(ctx, userInput)
}
