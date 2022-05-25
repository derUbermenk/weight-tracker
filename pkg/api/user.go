package api

import (
	"errors"
	"strings"
)

// UserService contains the methods of the user service
type UserService interface {
	New(user NewUserRequest) (createdUserID int, err error)
	Update(user UpdateUserRequest) (User, error)
	GetUser(id int) (user User, err error)
	All() (users []User, err error)
}

// UserRepository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface {
	CreateUser(NewUserRequest) (userID int, err error)
	UpdateUser(UpdateUserRequest) (User, error)
	GetUser(userID int) (User, error)
	GetUserByEmail(userEmail string) (user User, err error)
	GetUsers() ([]User, error)
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

func (u *userService) New(user NewUserRequest) (createdUserID int, err error) {
	// do some basic validations
	if user.Email == "" {
		err = errors.New("user service - email required")
		return
	}

	if user.Name == "" {
		err = errors.New("user service - name required")
		return
	}

	if user.WeightGoal == "" {
		err = errors.New("user service - weight goal required")
		return
	}

	var exists bool
	exists, err = emailExists(u.storage.GetUserByEmail, user.Email)

	if err != nil {
		return
	} else if exists {
		err = errors.New("user service - user with email already exists")
		return
	}

	// do some basic normalisation
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	createdUserID, err = u.storage.CreateUser(user)

	if err != nil {
		return
	}

	return
}

type userGetterByEmail func(email string) (user User, err error)

func emailExists(fn userGetterByEmail, email string) (exists bool, err error) {
	var user User
	user, err = fn(email)

	if err != nil {
		return
	}

	// proceed with comparison
	exists = user != User{}

	return
}
