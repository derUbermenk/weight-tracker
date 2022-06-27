package api_test

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"weight-tracker/pkg/api"
)

type mockUserRepo struct {
	users map[int]api.User
}

var taken_email = "taken_email@email.com"

var users = map[int]api.User{
	1: {
		ID:            1,
		Name:          "Rabbit",
		Age:           2,
		Height:        3,
		Sex:           "female",
		ActivityLevel: 2,
		WeightGoal:    "heavy",
		Email:         "some_email@email.com",
	},
	2: {
		ID:            2,
		Name:          "Mole",
		Age:           2,
		Height:        3,
		Sex:           "female",
		ActivityLevel: 2,
		WeightGoal:    "heavy",
		Email:         taken_email,
	},
}

func (m mockUserRepo) CreateUser(request api.NewUserRequest) (userID int, err error) {
	return userID, nil
}

func (m mockUserRepo) CreateUser_v2(email, hashedPassword string) (user api.User, err error) {
	user.Email = email
	user.HashedPassword = hashedPassword

	return
}

func (m mockUserRepo) GetUser(userID int) (api.User, error) {
	return m.users[userID], nil
}

func (m mockUserRepo) GetUserByEmail(userEmail string) (api.User, error) {
	// iterate over the items in m.users
	// check email, and return email if theirs
	for _, user := range m.users {
		if user.Email == userEmail {
			fmt.Printf("is user: %v \n\tgiven: %v\n\tcurrent: %v\n", userEmail == user.Email, userEmail, user.Email)
			return user, nil
		}
	}

	return api.User{}, nil
	/*
		if userEmail == taken_email {
			return m.users[2], nil // 2 is assigned the taken email
		}
	*/
}

func (m mockUserRepo) UpdateUser(request api.UpdateUserRequest) (api.User, error) {
	// assuming update has been validated
	// create the new user struct and make it the value
	// of the key identified by the user request key
	user_update := api.User{
		ID: request.ID, Name: request.Name,
		Age: request.Age, Height: request.Height,
		Sex: request.Sex, ActivityLevel: request.ActivityLevel,
		Email: request.Email, WeightGoal: request.WeightGoal,
	}
	m.users[request.ID] = user_update

	return m.users[request.ID], nil
}

func (m mockUserRepo) GetUsers() (users []api.User, err error) {
	// iterate over m.users map, and add all the values to the returned
	// users slice

	/* this does not work since mapping is unordered
	for _, user := range m.users {
		users = append(users, user)
	}
	*/

	// instead i iterate using the lenght of the users map
	for i := 1; i <= len(m.users); i++ {
		users = append(users, m.users[i])
	}

	return
}

func (m mockUserRepo) DeleteUser(userID int) (deletedUserID int, err error) {
	_, present := m.users[userID]

	if !present {
		return 0, nil
	}

	return userID, nil
}

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

	mockRepo := mockUserRepo{}
	userService := api.NewUserService(&mockRepo)
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
	userService := api.NewUserService(&userRepo)
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
	tests := []struct {
		name           string
		userEmail      string
		want_existence bool
		want_error     error
	}{
		{
			name:           "should return true when user exists",
			userEmail:      "taken_email@email.com",
			want_existence: true,
			want_error:     nil,
		},
		{
			name:           "should return false when user does not exist",
			userEmail:      "non_existent_user@email.com",
			want_existence: false,
			want_error:     nil,
		},
	}

	userRepo := mockUserRepo{users: users}
	userService := api.NewUserService(&userRepo)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			existence, err := userService.UserExists(test.userEmail)
			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if existence != test.want_existence {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, existence, test.want_existence)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	// check if the given password does not pass website bare minimums
	tests := []struct {
		name       string
		password   string
		want_valid bool
		want_error error
	}{
		{
			name:       "Should return false for a password less than 6 characters in len",
			password:   "lofi1",
			want_valid: false,
		},
		{
			name:       "Should return true for a password more than 6 characters",
			password:   "asdf234",
			want_valid: true,
		},
	}

	userService := api.NewUserService(nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid := userService.ValidatePassword(test.password)

			if valid != test.want_valid {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, valid, test.want_valid)
			}
		})
	}
}

func TestGetAllUser(t *testing.T) {
	tests := []struct {
		name       string
		request    api.NewUserRequest
		want_users []api.User
		want_error error
	}{
		{
			name:       "should return no users when there are none",
			want_users: nil, // changed from []api.User{} to nil because https://stackoverflow.com/questions/64643402/reflect-deepequal-is-returning-false-but-slices-are-same
			// getUsers mockRepo method returns nil slice when no user is found
			want_error: nil,
		}, {
			name: "should return users when there are users",
			want_users: []api.User{
				users[1],
				users[2],
			},
			want_error: nil,
		},
	}

	for _, test := range tests {
		test_users := copyUserMap(users)
		mockRepo := mockUserRepo{users: test_users}
		mockUserService := api.NewUserService(&mockRepo)

		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "should return no users when there are none":
				mockRepo.users = map[int]api.User{} // make sure there are no users
			case "should return users when there are users":
				mockRepo.users = test_users // use the predefined users
			}

			queried_users, err := mockUserService.All()

			if !reflect.DeepEqual(err, test.want_error) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_error)
			}

			if !reflect.DeepEqual(queried_users, test.want_users) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, queried_users, test.want_users)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		request    api.UpdateUserRequest
		want_user  api.User
		want_error error
	}{
		{
			// updates used to not work when the user decides not to change own email
			name: "should not return an error when the email is unchanged",
			request: api.UpdateUserRequest{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "some_email@email.com",
			},
			want_user: api.User{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "some_email@email.com",
			},
			want_error: nil,
		},
		{
			// updates used to not work when the user decides not to change own email
			name: "should update user when there is no conflict in email change",
			request: api.UpdateUserRequest{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "non_conflicting@email.com",
			},
			want_user: api.User{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "non_conflicting@email.com",
			},
			want_error: nil,
		},
		{
			name: "should not return an error when the email is changed but it does not exist yet",
			request: api.UpdateUserRequest{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "unused@email.com",
			},
			want_user: api.User{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         "unused@email.com",
			},
			want_error: nil,
		},
		{
			name: "should return an error because there is an existing email",
			request: api.UpdateUserRequest{
				ID:            1,
				Name:          "rabbit",
				Age:           20,
				Height:        250,
				Sex:           "male",
				WeightGoal:    "maintain",
				ActivityLevel: 2,
				Email:         taken_email,
			},
			want_user:  api.User{},
			want_error: errors.New("user service - user with email already exists"),
		},
	}

	for _, test := range tests {
		test_users := copyUserMap(users)
		mockRepo := mockUserRepo{users: test_users}
		mockUserService := api.NewUserService(&mockRepo)

		t.Run(test.name, func(t *testing.T) {
			user, err := mockUserService.Update(test.request)

			if !reflect.DeepEqual(err, test.want_error) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_error)
			}

			if !reflect.DeepEqual(user, test.want_user) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, user, test.want_user)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {

	tests := []struct {
		name       string
		request    int
		want_error error
		want_id    int
	}{
		{
			name:       "should delete user successfully since user exists",
			request:    1,
			want_error: nil,
			want_id:    1,
		},
		{
			name:       "should return an error when user with submitted id does not exist",
			request:    25,
			want_error: errors.New("user service - user with given id does not exist"),
			want_id:    0,
		},
	}

	for _, test := range tests {
		test_users := copyUserMap(users)
		mockRepo := mockUserRepo{users: test_users}
		mockUserService := api.NewUserService(&mockRepo)

		t.Run(test.name, func(t *testing.T) {
			userID, err := mockUserService.Delete(test.request)

			if !reflect.DeepEqual(err, test.want_error) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_error)
			}

			if !reflect.DeepEqual(userID, test.want_id) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, userID, test.want_id)
			}
		})
	}
}

// convenience function for copying user map therefore isolating changes to tests
func copyUserMap(source_map map[int]api.User) (copied_map map[int]api.User) {
	copied_map = make(map[int]api.User)

	for k, v := range source_map {
		copied_map[k] = v
	}

	return
}
