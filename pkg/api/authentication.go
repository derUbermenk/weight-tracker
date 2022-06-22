package api

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenClaims struct {
	TokenType string `json:"token_type"`
	Email     string `json:"email"`

	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	Email     string `json:"email"`
	CustomKey string `json:"custom_key"`

	// this field is required for it to be considered a RefreshTokenClaims
	// also allows us additional fields
	jwt.RegisteredClaims
}

var TokenStatus = map[string]int{
	"valid":   0,
	"invalid": 1,
	"expired": 2,
}

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
	GenerateAccessToken(email string, expiration int64) (signed_access_token string, err error)
	GenerateRefreshToken(email string, customKey string) (signed_refresh_token string, err error)
	ValidateAccessToken(access_token string) (status int)
	ValidateRefreshToken(refresh_token, custom_key string) (validity bool)
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

// generates an access token and adds the email and expiration to the claims(payload of the token).
// expiration is an int which represents unix seconds since 1970, January 1.
func (a *authService) GenerateAccessToken(email string, expiration int64) (signed_access_token string, err error) {
	claims := &AccessTokenClaims{
		TokenType: "access",
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiration, 0)),
		},
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_access_token, err = access_token.SignedString(a.signingKey_byte)

	if err != nil {
		log.Printf("Service Error: %v", err)
		return
	}

	return
}

func (a *authService) GenerateRefreshToken(email string, customKey string) (signed_refresh_token string, err error) {
	claims := &RefreshTokenClaims{
		TokenType: "refresh",
		Email:     email,
		CustomKey: customKey,
	}

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_refresh_token, err = refresh_token.SignedString(a.signingKey_byte)

	if err != nil {
		log.Printf("Service Error: %v", err)
		return
	}

	return
}

func (a *authService) ValidateAccessToken(access_token string) (status int) {
	claims := &AccessTokenClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return a.signingKey_byte, nil
	}

	_, err := jwt.ParseWithClaims(access_token, claims, keyFunc)

	if err != nil {
		log.Printf("Service Error: %v", err)

		// err could also be caused when the jwt token is expired
		// but compare the error with jwt.ErrTokenExpired does not work
		// to augment this, I resolved to check the error value instead
		if errors.Is(err, jwt.ErrTokenExpired) {
			status = TokenStatus["expired"]
		} else {
			status = TokenStatus["tampered"]
		}

		return
	}

	status = TokenStatus["valid"]
	return
}

func (a *authService) ValidateRefreshToken(refresh_token, customKey string) (validity bool) {
	claims := &RefreshTokenClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return a.signingKey_byte, nil
	}

	tkn, err := jwt.ParseWithClaims(refresh_token, claims, keyFunc)

	if err != nil {
		log.Printf("service error: %v", err)
		return
	} else if !tkn.Valid {
		return
	}

	if claims.CustomKey != customKey {
		return
	}

	validity = true
	return
}

// compares the the password with the hashedPassword
func (a *authService) authenticate(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}
