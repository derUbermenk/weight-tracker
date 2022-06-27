package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// UserService contains the methods of the user service
type UserService interface {
	New(user NewUserRequest) (createdUserID int, err error)
	Delete(userID int) (deletedUserID int, err error)
	Update(user UpdateUserRequest) (User, error)
	GetUser(id int) (user User, err error)
	All() (users []User, err error)

	CreateUser(email, hashedPassword string) (user User, err error)
	HashPassword(password string) (hashedPass string, err error)
	UserExists(email string) (exists bool, err error)
	ValidatePassword(password string) (validity bool)
}

// UserRepository is what lets our service do db operations without knowing anything about the implementation
type UserRepository interface {
	CreateUser(NewUserRequest) (userID int, err error)
	DeleteUser(userID int) (deletedUserID int, err error)
	UpdateUser(UpdateUserRequest) (User, error)
	GetUser(userID int) (User, error)
	GetUserByEmail(userEmail string) (user User, err error)
	GetUsers() ([]User, error)

	CreateUser_v2(email, hashedPassword string) (User, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) CreateUser(email, hashedPassword string) (user User, err error) {
	user, err = u.storage.CreateUser_v2(email, hashedPassword)

	if err != nil {
		log.Printf("Service error: %v", err)
		return
	}

	return
}

func (u *userService) Update(user UpdateUserRequest) (updatedUser User, err error) {
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	var exists bool
	var changed bool

	changed, err = emailChanged(u.storage.GetUser, user.ID, user.Email)
	exists, err = emailExists(u.storage.GetUserByEmail, user.Email)

	if err != nil {
		return
	} else if changed && exists {
		err = errors.New("user service - user with email already exists")
		return
	}

	updatedUser, err = u.storage.UpdateUser(user)

	if err != nil {
		return
	}

	return
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

func (u *userService) Delete(userID int) (deletedUserID int, err error) {
	deletedUserID, err = u.storage.DeleteUser(userID)

	if err != nil {
		return
	} else if deletedUserID == 0 {
		err = errors.New("user service - user with given id does not exist")
		return
	}

	return
}

func (u *userService) HashPassword(password string) (hashedPass string, err error) {
	var hashedPass_byte []byte
	hashedPass_byte, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Service error: %v")
		return
	}

	hashedPass = string(hashedPass_byte)
	return
}

func (u *userService) UserExists(email string) (exists bool, err error) {
	var user User
	user, err = u.storage.GetUserByEmail(email)

	if err != nil {
		log.Printf("Service Error: %v", err)
		return
	}

	fmt.Printf("exists: %v\n", user != User{})
	fmt.Printf("%v \n*************\n", user)

	exists = user != User{}
	return
}

func (u *userService) ValidatePassword(password string) (validity bool) {
	if len(password) >= 6 {
		validity = true
		return
	}

	return
}

type userGetterByEmail func(email string) (user User, err error)

// checks if the email submitted is already used
func emailExists(userGetter userGetterByEmail, email string) (exists bool, err error) {
	var user User
	user, err = userGetter(email)

	if err != nil {
		return
	}

	// proceed with comparison
	exists = user != User{}

	return
}

type userGetter func(id int) (user User, err error)

// checks if the submitted email is not the same as the users current email
func emailChanged(userGetter userGetter, requestID int, requestEmail string) (unchanged bool, err error) {
	var user User
	user, err = userGetter(requestID) // get user

	if err != nil {
		return
	}

	unchanged = requestEmail != user.Email

	return
}
