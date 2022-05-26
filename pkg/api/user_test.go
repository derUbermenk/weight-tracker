package api_test

import (
	"errors"
	"reflect"
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

func (m mockUserRepo) GetUser(userID int) (api.User, error) {
	return m.users[userID], nil
}

func (m mockUserRepo) GetUserByEmail(userEmail string) (api.User, error) {
	// iterate over the items in m.users
	// check email, and return email if theirs
	for _, user := range m.users {
		if user.Email == userEmail {
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

func TestCreateNewUser(t *testing.T) {

	tests := []struct {
		name     string
		request  api.NewUserRequest
		want_err error
		want_id  int
	}{
		{
			name: "should create a new user successfully",
			request: api.NewUserRequest{
				Name:          "test user",
				WeightGoal:    "maintain",
				Age:           20,
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want_err: nil,
			want_id:  0,
		}, {
			name: "should return an error because of missing email",
			request: api.NewUserRequest{
				Name:          "test user",
				Age:           20,
				WeightGoal:    "maintain",
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "",
			},
			want_err: errors.New("user service - email required"),
			want_id:  0,
		}, {
			name: "should return an error because of missing name",
			request: api.NewUserRequest{
				Name:          "",
				Age:           20,
				WeightGoal:    "maintain",
				Height:        180,
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want_err: errors.New("user service - name required"),
			want_id:  0,
		}, {
			name: "should return error because user with email already exists",
			request: api.NewUserRequest{
				Name:          "test user with email exists",
				Age:           20,
				Height:        180,
				WeightGoal:    "maintain",
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "taken_email@email.com",
			},
			want_err: errors.New("user service - user with email already exists"),
			want_id:  0,
		},
	}

	for _, test := range tests {
		test_users := copyUserMap(users)
		mockRepo := mockUserRepo{users: test_users}
		mockUserService := api.NewUserService(&mockRepo)

		t.Run(test.name, func(t *testing.T) {
			userID, err := mockUserService.New(test.request)

			if !reflect.DeepEqual(err, test.want_err) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_err)
			}

			if userID != test.want_id {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_err)
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

// convenience function for copying user map therefore isolating changes to tests
func copyUserMap(source_map map[int]api.User) (copied_map map[int]api.User) {
	copied_map = make(map[int]api.User)

	for k, v := range source_map {
		copied_map[k] = v
	}

	return
}
