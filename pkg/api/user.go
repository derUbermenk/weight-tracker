package api

import (
	"errors"
	"strings"
)

// UserService contains the methods of the user service
type UserService interface {
	New(user NewUserRequest) error
	Update(user UpdateUserRequest) error
	GetUser(id int) (user User, err error)
	All() (users []User, err error)
}

// UserRepository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface {
	CreateUser(NewUserRequest) error
	UpdateUser(UpdateUserRequest) error
	GetUser(userID int) (User, error)
	GetUsers() ([]User, error)
	GetUserByEmail(userEmail string) (User, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) Update(user UpdateUserRequest) error {
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	err := u.storage.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userService) GetUser(userID int) (User, error) {
	user, err := u.storage.GetUser(userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *userService) All() ([]User, error) {
	users, err := u.storage.GetUsers()

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (u *userService) New(user NewUserRequest) error {
	// do some basic validations
	if user.Email == "" {
		return errors.New("user service - email required")
	}

	if user.Name == "" {
		return errors.New("user service - name required")
	}

	if user.WeightGoal == "" {
		return errors.New("user service - weight goal required")
	}

	existingEmail, err := emailExists(user.Email, u.storage)

	if err != nil {
		return err
	}

	if existingEmail {
		return errors.New("user service - user with email already exists")
	}

	// do some basic normalisation
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	err = u.storage.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func getUserByEmail(email string, storage UserRepository) (User, error) {
	user, err := storage.GetUserByEmail(email)

	return user, err
}

func emailExists(email string, storage UserRepository) (bool, error) {
	user, err := storage.GetUserByEmail(email)

	return user != User{}, err
}
