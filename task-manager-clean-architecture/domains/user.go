package domains

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Role     string             `bson:"role" json:"role"`
}

// NewUser creates a new User instance with hashed password
func NewUser(username, password, role string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}, nil
}

// hashPassword hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	UserExists(ctx context.Context, username string) (bool, error)
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, username, password, role string) (*User, error)
	AuthenticateUser(ctx context.Context, username, password string) (*User, bool)
}
