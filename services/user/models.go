package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role string

const (
	RoleRoot    Role = "ROOT"
	RoleDoctor  Role = "DOCTOR"
	RoleStaff   Role = "STAFF"
	RolePatient Role = "PATIENT"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	ProfilePic string             `bson:"profilePic" json:"profilePic"`
	Password   string             `bson:"password" json:"-"`
	Role       Role               `bson:"role" json:"role"`
}

type Doctor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Specialty string             `bson:"specialty" json:"specialty"`
	Documents []string           `bson:"documents" json:"documents"`
}

type UserInput struct {
	Name       string
	Email      string
	ProfilePic string
	Password   string
	Role       Role
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token string
	User  PublicUser
}

// PublicUser DTO without sensitive fields
type PublicUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Role       Role   `json:"role"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:         u.ID.Hex(),
		Name:       u.Name,
		Email:      u.Email,
		ProfilePic: u.ProfilePic,
		Role:       u.Role,
	}
}
