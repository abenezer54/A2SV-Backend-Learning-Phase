package service

import (
	"context"
	"errors"

	"task-manager-api/models"
	"task-manager-api/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, username, password, role string) (*models.User, error) {
	// Check if the username already exists
	exists, err := s.userRepository.UserExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Create a new user
	user, err := models.NewUser(username, password, role)
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

func (s *UserService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, bool) {
	user, err := s.userRepository.FindUserByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, false
	}
	if user.CheckPassword(password) {
		return user, true
	}
	return nil, false
}
