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
			customKey:          "hashed_email_and_hashed_pass",
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
	// the service setup code was placed here because one of the tokens require the service
	// to generate valid access tokens
	mockRepo := mockAuthRepo{}
	authService := api.NewAuthService(&mockRepo, signingKey)

	// created by generating expired tokens in jwt.io. This still used the signing key used in the test.
	expired_tokens := [3]string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJleGlzdGluZ19lbWFpbEBlbWFpbC5jb20iLCJleHAiOjE2NTU2NTgwfQ.I_bv7NhYuRZSSuIjlFHlA6fRqBuXXblY1a28AVYXEZ0",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJuZXdfZW1haWxAZW1haWwuY29tIiwiZXhwIjoyOTk5ODkxfQ.KEq20X7eUBidx7avrFce4Jum36gy4j6johqsh37ggzc",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJhbm90aGVyX2VtYWlsQGVtYWlsLmNvbSIsImV4cCI6Mjk4OTExMjMxfQ.o9pC7sH4wwBLaOgq4q8p4z69Yaq_PxI6ki5MwDyfAYQ",
	}

	// created by generating a valid token using jwt.io, sharing the token, then changing the payload.
	// this simulates using the token to login as another user.
	tampered_tokens := [3]string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJteV9lbWFpbEBlbWFpbC5jb20iLCJleHAiOjI5ODI5MTEyMzF9.FS47rXB1qvb81wg0h4EIu8MSvQefI0-LHyr3vz0NwC4",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiJ0YW1wZXJlZF9lbWFpbEBlbWFpbC5jb20iLCJleHAiOjI5ODI5MTEyMzF9.Zf6wzs_pZxJpOOicMOEgRSDMfNEvPrtOuqllHIVO9ZA",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZW1haWwiOiIxMXRhbXBlcmVkX2VtYWlsQGVtYWlsLmNvbSIsImV4cCI6Mjk4MjkxMTIzMX0.USV6CNMJy73vlahpj79Bh-Hy3jLfjrfo5e_o_8jlSF8",
	}

	// created by using the GenerateAccessToken method of the auth api. This is also a good test to check wether
	// the generated access tokens are indeed valid
	var valid_tokens [3]string
	emails := [3]string{"newEmail@io.com", "newWark@email.com", "noTraffic@email.com"}
	expirations := [3]int64{2982911231, 3982511231, 3782511231}

	for i := 0; i <= 2; i++ {
		access_token, _ := authService.GenerateAccessToken(emails[i], expirations[i])
		valid_tokens[i] = access_token
	}

	// define the tests

	tests := []struct {
		name          string
		access_tokens [3]string
		want_status   int
		want_err      error
	}{
		{
			name:          "Must all have expired status",
			access_tokens: expired_tokens,
			want_status:   api.TokenStatus["expired"],
			// want_err:      nil,
		},
		{
			name:          "Must all have tampered status",
			access_tokens: tampered_tokens,
			want_status:   api.TokenStatus["tampered"],
			// want_err:      nil,
		},
		{
			name:          "Must all have valid status",
			access_tokens: valid_tokens,
			want_status:   api.TokenStatus["valid"],
			// want_err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for index, token := range test.access_tokens {
				status := authService.ValidateAccessToken(token)

				// remove error returns for now
				/*
					if err != test.want_err {
						t.Errorf("test %v failed.\n\tgot %v,\t\nwanted: %v\n\tOn token: %v", test.name, err, test.want_err, index)
					}
				*/

				if status != test.want_status {
					t.Errorf("test %v failed.\n\tgot %v,\t\nwanted: %v\n\tOn token: %v", test.name, status, test.want_status, index)
				}
			}
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {

	// the custom key is the combination of the users email and hashed password
	// we do not generate this in this case, and only use a random string
	// what's important for refresh token validity is that the custom key is the same
	// as that in the payload.
	var correct_custom_key string
	var incorrect_custom_key string

	// the following variables each represent a set of tokens with a given condition
	//
	// valid_tokens: the tokens where generated using the signing key and the custom key used is unchanged
	// wrongly_signed_tokens: the token was not generated using the signing key
	// tampered_tokens: the signed_tokens strings have been changed.
	var valid_tokens [3]string
	var wrongly_signed_tokens [3]string
	var tampered_tokens [3]string

	correct_custom_key = "correct_custom_key"
	incorrect_custom_key = "incorrect_custom_key"

	// all these tokens used the same custom key but with different emails. The signing key used is the singning_key global variable available in this test.
	// follow the links to see the decoded version
	valid_tokens = [3]string{
		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdfZW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.ThLgaHwynarnJYSrLT360Ki3JMrD9Jgb_DL1C9NN3HM
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdfZW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.ThLgaHwynarnJYSrLT360Ki3JMrD9Jgb_DL1C9NN3HM",

		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3RW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.GEj4aXmKc1JLDKE1xIBsZoYyH7PggEavnEzDsFzVn5s
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3RW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.GEj4aXmKc1JLDKE1xIBsZoYyH7PggEavnEzDsFzVn5s",

		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3MkVtYWlsQGVtYWlsLmNvbSIsImN1c3RvbV9rZXkiOiJjb3JyZWN0X2N1c3RvbV9rZXkifQ.DSoTvVEaPEU1Mrz2O2gGp7MTShkMU5I2KIhoNme2BIY
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3MkVtYWlsQGVtYWlsLmNvbSIsImN1c3RvbV9rZXkiOiJjb3JyZWN0X2N1c3RvbV9rZXkifQ.DSoTvVEaPEU1Mrz2O2gGp7MTShkMU5I2KIhoNme2BIY",
	}

	// these tokens are generated using the jwt.io websites decode and encode function with different signing keys
	wrongly_signed_tokens = [3]string{
		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3MkVtYWlsQGVtYWlsLmNvbSIsImN1c3RvbV9rZXkiOiJjb3JyZWN0X2N1c3RvbV9rZXkifQ.0EWQpapPJFhwb47JooN2d7NpueXhQ8ZMyvsbrOnpCFQ
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3MkVtYWlsQGVtYWlsLmNvbSIsImN1c3RvbV9rZXkiOiJjb3JyZWN0X2N1c3RvbV9rZXkifQ.0EWQpapPJFhwb47JooN2d7NpueXhQ8ZMyvsbrOnpCFQ",

		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdFbWFpbEBlbWFpbC5jb20iLCJjdXN0b21fa2V5IjoiY29ycmVjdF9jdXN0b21fa2V5In0.hd5Rym36TeMYrkJJZb4KvsN2NM3jVTkdQ04mP-LpunA
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdFbWFpbEBlbWFpbC5jb20iLCJjdXN0b21fa2V5IjoiY29ycmVjdF9jdXN0b21fa2V5In0.hd5Rym36TeMYrkJJZb4KvsN2NM3jVTkdQ04mP-LpunA",

		// https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdFbWFpbEBlbWFpbC5jb20iLCJjdXN0b21fa2V5IjoiY29ycmVjdF9jdXN0b21fa2V5In0.tgGXUh6302JdUFCYB4gp_yBSbHYghag6TzIS_XoSWk8
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdFbWFpbEBlbWFpbC5jb20iLCJjdXN0b21fa2V5IjoiY29ycmVjdF9jdXN0b21fa2V5In0.tgGXUh6302JdUFCYB4gp_yBSbHYghag6TzIS_XoSWk8",
	}

	// this tokens are from the valid tokens but with with randomly additional or removed characters.
	tampered_tokens = [3]string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9_.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoiZXhpc3RpbmdfZW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.ThLgaHwynarnJYSrLT360Ki3JMrD9Jgb_DL1C9NN3HM",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eXJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3RW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.GEj4aXmKc1JLDKE1xIBsZoYyH7PggEavnEzDsFzVn5s",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eXJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImVtYWlsIjoibmV3RW1haWxAZW1haWwuY29tIiwiY3VzdG9tX2tleSI6ImNvcnJlY3RfY3VzdG9tX2tleSJ9.GEj4aXmKc1JLDKE1xIBsZoYyH7PggEavnEzDsFzVn5!s",
	}

	tests := []struct {
		name          string
		custom_key    string
		tokens        [3]string
		want_validity bool
		want_err      error
	}{
		{
			name:          "Must be valid for valid tokens with the same custom keys",
			custom_key:    correct_custom_key,
			tokens:        valid_tokens,
			want_validity: true,
			want_err:      nil,
		},
		{
			name:          "Must be invalid for valid tokens with different custom keys",
			custom_key:    incorrect_custom_key,
			tokens:        valid_tokens,
			want_validity: false,
			want_err:      nil,
		},
		{
			name:          "Must be invalid for wrongly signed tokens",
			custom_key:    correct_custom_key,
			tokens:        wrongly_signed_tokens,
			want_validity: false,
			want_err:      nil,
		},
		{
			name:          "Must be invalid for tampered tokens",
			custom_key:    correct_custom_key,
			tokens:        tampered_tokens,
			want_validity: false,
			want_err:      nil,
		},
	}

	// setup outside of the loop to save us a few steps
	// only add this inside when we are dealing with database manipulation
	// methods
	mockRepo := mockAuthRepo{}
	authService := api.NewAuthService(&mockRepo, signingKey)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for index, token := range test.tokens {
				// check the validity and error
				validity := authService.ValidateRefreshToken(token, test.custom_key)

				/*
					if err != test.want_err {
						t.Errorf("test %v failed.\n\tgot %v,\t\nwanted: %v\n\tOn token: %v", test.name, err, test.want_err, index)
					}
				*/

				if validity != test.want_validity {
					t.Errorf("test %v failed.\n\tgot %v,\t\nwanted: %v\n\tOn token: %v", test.name, validity, test.want_validity, index)
				}
			}
		})
	}
}
