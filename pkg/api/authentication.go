package api

import (
	"golang.org/x/crypto/bcrypt"
)

// authentication interface. Defines the
// 	functions and its signature that a type needs
// 	to implement to be able to be considered
// 	as such
//
//	the Authentication service handles user
//  authentication api logic such as jwt access
// 	and refresh token validation and creation.
// 	As well as Credential(password and email) validation.
type AuthService interface {
	ValidateCredentials(email, password string) (validity bool, err error)
}

// authentication repository interface represents any
//	type that is used by the authentication service to
// 	interact with the system db.
type AuthRepository interface {
	GetUserByEmail(email string) (user User, err error)
}

// a struct type representing the authentication
// service interface
type authService struct {
	storage         AuthRepository
	signingKey_byte []byte
}

// creates a new authservice for use of the server
func NewAuthService(authRepository AuthRepository, signingKey string) AuthService {
	return &authService{
		storage:         authRepository,
		signingKey_byte: []byte(signingKey),
	}
}

func (a *authService) ValidateCredentials(email, password string) (validity bool, err error) {
	// find user email in db
	user, err := a.storage.GetUserByEmail(email)

	// if err exists or user with email does not exist
	// return
	// note that validity zero value is false
	if err != nil || (user == User{}) {
		return
	}

	// otherwise
	// compare user.password with password
	// if not equal set validity == false
	validity = a.authenticate(user.Password, password)

	return
}

// compares the the password with the hashedPassword
func (a *authService) authenticate(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}
