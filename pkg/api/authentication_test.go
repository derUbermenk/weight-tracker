package api_test

import (
	"fmt"
	"testing"
	"weight-tracker/pkg/api"

	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "super_secret_hs256"
)

type Credentials struct {
	Email    string
	Password string
}

type mockAuthRepo struct {
	users map[string]string
}

func (m *mockAuthRepo) GetUserByEmail(email string) (user api.User, err error) {
	password, exists := m.users[email]

	if !exists {
		return api.User{}, nil
	}

	user.Email = email
	user.Password = password

	return
}

func NewMockAuthRepo() (authRepo *mockAuthRepo) {
	authRepo = &mockAuthRepo{
		users: map[string]string{},
	}

	users_unstored := map[string]string{
		"existing_email@email.com": "correct_password1234",
	}

	// hash the passwords
	for email, password := range users_unstored {
		hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		if err != nil {
			fmt.Print("here")
			panic(err)
		}

		authRepo.users[email] = string(hashed_pass)
	}

	return
}

func TestValidateCredentials(t *testing.T) {
	// compares the credentials passed with the credentials in store

	tests := []struct {
		name          string
		credentials   Credentials
		want_validity bool
		want_err      error
	}{
		{
			name:          "should return true if the credentials passed match those of store",
			credentials:   Credentials{Email: "existing_email@email.com", Password: "correct_password1234"},
			want_validity: true,
			want_err:      nil,
		},
		{
			name:          "should return false if the credential password does not match does of db",
			credentials:   Credentials{Email: "existing_email@email.com", Password: "incorrect_password1234"},
			want_validity: false,
			want_err:      nil,
		},
		{
			name:          "should return false if the credential email does not exist",
			credentials:   Credentials{Email: "nonexisting_email@email.com", Password: "incorrect_password1234"},
			want_validity: false,
		},
	}

	for _, test := range tests {
		mockRepo := NewMockAuthRepo()
		authService := api.NewAuthService(mockRepo, signingKey)

		t.Run(test.name, func(t *testing.T) {
			validity, err := authService.ValidateCredentials(test.credentials.Email, test.credentials.Password)

			if err != test.want_err {
				t.Errorf("test %v failed. got: %v, wanted: %v", test.name, err, test.want_err)
			}

			if validity != test.want_validity {
				t.Errorf("test %v failed. got: %v, wanted: %v", test.name, validity, test.want_validity)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	// generates an access token given the credentials
	tests := []struct {
		name              string
		email             string
		expiration        int64
		want_access_token string
		want_err          error
	}{
		{
			// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJleGlzdGluZ19lbWFpbEBlbWFpbC5jb20iLCJleHAiOjE2NTU2OTA1NzF9.KbC55juOe0dNX7DC7lT4vWaF-XhmNHXzi9UvqEQ1V1A
			name:              "should return the right token given the following credentials v1",
			email:             "existing_email@email.com",
			expiration:        1655690571,
			want_access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJleGlzdGluZ19lbWFpbEBlbWFpbC5jb20iLCJleHAiOjE2NTU2OTA1NzF9.KbC55juOe0dNX7DC7lT4vWaF-XhmNHXzi9UvqEQ1V1A",
			want_err:          nil,
		},
		{
			// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJuZXdFbWFpbEBlbWFpbC5pbyIsImV4cCI6MTAxMTkwOTAxfQ.b5IhQ6pig8fqUiKxNt3LqP0Cs_21pDiHLG4U32TVJDo
			name:              "should return the right token given the following credentials v2",
			email:             "newEmail@email.io",
			expiration:        101190901,
			want_access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJuZXdFbWFpbEBlbWFpbC5pbyIsImV4cCI6MTAxMTkwOTAxfQ.b5IhQ6pig8fqUiKxNt3LqP0Cs_21pDiHLG4U32TVJDo",
			want_err:          nil,
		},
	}

	for _, test := range tests {
		mockRepo := NewMockAuthRepo()
		authService := api.NewAuthService(mockRepo, signingKey)

		t.Run(test.name, func(t *testing.T) {
			generateAccessToken, err := authService.GenerateAccessToken(test.email, test.expiration)

			if err != test.want_err {
				t.Errorf("test %v failed. got %v, wanted: %v", test.name, err, test.want_err)
			}

			if generateAccessToken != test.want_access_token {
				t.Errorf("test %v failed.\n\tgot: %v\n\twanted: %v", test.name, generateAccessToken, test.want_access_token)
			}
		})
	}

}

func TestGenerateRefreshToken(t *testing.T) {
	// generates a refresh token given an email and a custom key
	tests := []struct {
		name               string
		email              string
		customKey          string
		want_refresh_token string
		want_err           error
	}{
		{
			name:               "Should return the correct refresh token given the parameters",
			email:              "existing_email@email.com",
			customKey:          "hashed_email_and_pass",
			want_refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdfZW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6Imhhc2hlZF9lbWFpbF9hbmRfaGFzaGVkX3Bhc3MifQ.q0WPazaOaGnnfTrsjAfApXEafYWdpNWwUHMMCQP7FB4",
			want_err:           nil,
		},
		{
			name:               "Should return the correct refresh token given the correct params",
			email:              "newEmail@email.io",
			customKey:          "hashed_userEmail_and_pass",
			want_refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3RW1haWxAZW1haWwuaW8iLCJjdXN0b21fa2V5IjoiaGFzaGVkX3VzZXJFbWFpbF9hbmRfcGFzcyJ9.CfwwZ-gq7rTti3OuWiPLUBAwmShhEK6K7p4bxo_zgsc",
			want_err:           nil,
		},
	}

	for _, test := range tests {
		mockRepo := NewMockAuthRepo()
		authService := api.NewAuthService(mockRepo, signingKey)

		t.Run(test.name, func(t *testing.T) {
			refresh_token, err := authService.GenerateRefreshToken(test.email, test.customKey)

			if err != test.want_err {
				t.Errorf("test %v failed.\n\tgot %v\n\twanted: %v", test.name, err, test.want_err)
			}

			if refresh_token != test.want_refresh_token {
				t.Errorf("test %v failed.\n\tgot %v,\n\twanted: %v", test.name, refresh_token, test.want_refresh_token)
			}
		})
	}

}

func TestValidateAccessToken(t *testing.T) {
	// validate access token

	//
}

func TestValidateRefreshToken(t *testing.T) {
	// validate access token

	//
}
