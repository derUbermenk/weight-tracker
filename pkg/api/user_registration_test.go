package api_test

import (
	"regexp"
	"testing"
	"weight-tracker/pkg/api"
)

/*
	this file contains test for user methods that are used in user registration.
	The methods tested are still methods of a UserService interface
*/

func TestCreateUser(t *testing.T) {
	// the CreateUser method handles the actual call to the database to save the user
	// on a successful creation it will not return an error otherwise, it will.

	// otherwise it returns an empty interface and an error

	// we do not expect any errors from the service itself, but it can return one assuming
	// the storage interface it uses encounters one.
	// for the tests however, it is

	// test 1
	tests := []struct {
		name           string
		email          string
		hashedPassword string
		want_user      api.User
		want_error     error
	}{
		{
			name:           "It creates and returns the user created",
			email:          "newEmail@email.com",
			hashedPassword: "xHaaPass",
			want_user:      api.User{Email: "newEmail@email.com", HashedPassword: "xHaaPass"},
			want_error:     nil,
		},
	}

	userRepo := mockUserRepo{}
	userService := api.NewUserService(userRepo)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := userService.CreateUser(test.email, test.hashedPassword)

			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if user != test.want_user {
				t.Errorf("test %v failed.\n\t got: %+v\n\twanted: %+v", test.name, user, test.want_user)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	// checks if the password is indeed hashed as intended.
	// not necessarily check if the hashed is the same as an intended hash.
	// just check if the return value follows some sort of pattern.

	tests := []struct {
		name            string
		password        string
		want_error      error
		want_hashLenght int

		// hashed passwords look like this $2a$10$DqbfQm0RQEGjvrPWv.IN.eNvk5VJ6g0A.DnN1g50jZS6L319l8GAC
		want_has_letters_numbers_and_symbols bool
	}{
		{
			name:                                 "It creates a hashed Password",
			password:                             "myPasswordRabbits",
			want_hashLenght:                      60,
			want_has_letters_numbers_and_symbols: true,
			want_error:                           nil,
		},
	}

	userRepo := mockUserRepo{}
	userService := api.NewUserService(userRepo)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// pattern represents letters numbers and symbols.
			// see this link for better explanation of the regex
			// hashing was successfull if there is a-zA-Z0-9_ or \/$.
			pattern, err := regexp.Compile(`\w[\/\\.\$]?`)

			if err != nil {
				t.Errorf("test %v failed.\n\tError: %v\n\tat regex compilation", test.name, err)
			}

			hashedPass, err := userService.HashPassword(test.password)
			hashLenght := len(hashedPass)
			has_letters_numbers_and_symbols := pattern.MatchString(hashedPass)

			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if hashLenght != test.want_hashLenght {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, hashLenght, test.want_has_letters_numbers_and_symbols)
			}

			if has_letters_numbers_and_symbols != test.want_has_letters_numbers_and_symbols {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, has_letters_numbers_and_symbols, test.want_has_letters_numbers_and_symbols)
			}

		})
	}
}

func TestUserExists(t *testing.T) {
	// checks if the function indeed returns a boolean
	// that validates the User's existence in storage
}
