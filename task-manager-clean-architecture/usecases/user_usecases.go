package usecases

import (
	"context"
	"errors"

	"task-manager-api/domains"
)

type UserUsecase struct {
	userRepository domains.UserRepository
}

func NewUserUsecase(userRepo domains.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepo,
	}
}

func (s *UserUsecase) RegisterUser(ctx context.Context, username, password, role string) (*domains.User, error) {
	// Check if the username already exists
	exists, err := s.userRepository.UserExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Create a new user
	user, err := domains.NewUser(username, password, role)
	if err != nil {
		return nil, err
	}

	// Save the new user to the database
	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserUsecase) AuthenticateUser(ctx context.Context, username, password string) (*domains.User, bool) {
	user, err := s.userRepository.FindUserByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, false
	}
	if user.CheckPassword(password) {
		return user, true
	}
	return nil, false
}
