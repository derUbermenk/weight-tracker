package api_test

import (
	"errors"
	"reflect"
	"testing"
	"weight-tracker/pkg/api"
)

type mockUserRepo struct {
	users []api.User
}

var taken_email = "taken_email@email.com"

var users = []api.User{
	api.User{
		ID:            1,
		Name:          "Rabbit",
		Age:           2,
		Height:        3,
		Sex:           "female",
		ActivityLevel: 2,
		WeightGoal:    "heavy",
		Email:         "some_email@email.com",
	},
	api.User{
		ID:            1,
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
	return api.User{}, nil
}

func (m mockUserRepo) GetUserByEmail(userEmail string) (api.User, error) {
	if userEmail == taken_email {
		return m.users[1], nil
	}

	return api.User{}, nil
}

func (m mockUserRepo) UpdateUser(request api.UpdateUserRequest) (api.User, error) {
	return api.User{}, nil
}

func (m mockUserRepo) GetUsers() (users []api.User, err error) {
	users = m.users
	return
}

func TestCreateNewUser(t *testing.T) {
	mockRepo := mockUserRepo{users: users}
	mockUserService := api.NewUserService(&mockRepo)

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
	mockRepo := mockUserRepo{}
	mockUserService := api.NewUserService(&mockRepo)

	tests := []struct {
		name       string
		request    api.NewUserRequest
		want_users []api.User
		want_error error
	}{
		{
			name:       "should return no users when there are none",
			want_users: []api.User{},
			want_error: nil,
		}, {
			name:       "should return users when there are users",
			want_users: users,
			want_error: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "should return no users when there are none":
				mockRepo.users = []api.User{} // make sure there are no users
			case "should return users when there are users":
				mockRepo.users = users // use the predefined users
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
	mockRepo := mockUserRepo{users: users}
	mockUserService := api.NewUserService(&mockRepo)

	tests := []struct {
		name       string
		request    api.UpdateUserRequest
		want_user  api.User
		want_error error
	}{
		{
			name: "should return an error because there is an existing email",
			request: api.UpdateUserRequest{
				ID:         1,
				Name:       "rabbit",
				Age:        20,
				Height:     250,
				Sex:        "male",
				WeightGoal: "maintain",
				Email:      taken_email,
			},
			want_user:  api.User{},
			want_error: errors.New("user service - user with email already exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := mockUserService.Update(test.request)

			if errors.Is(err, test.want_error) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want_error)
			}

			if user != test.want_user {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, user, test.want_user)
			}
		})
	}
}
