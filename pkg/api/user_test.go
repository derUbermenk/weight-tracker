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

func (m mockUserRepo) CreateUser(request api.NewUserRequest) error {
	if request.Name == "test user already created" {
		return errors.New("repository - user already exists in database")
	}

	return nil
}

func (m mockUserRepo) GetUsers() (users []api.User, err error) {
	users = m.users
	return
}

func TestCreateNewUser(t *testing.T) {
	mockRepo := mockUserRepo{}
	mockUserService := api.NewUserService(&mockRepo)

	tests := []struct {
		name    string
		request api.NewUserRequest
		want    error
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
			want: nil,
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
			want: errors.New("user service - email required"),
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
			want: errors.New("user service - name required"),
		}, {
			name: "should return error from database because user already exists",
			request: api.NewUserRequest{
				Name:          "test user already created",
				Age:           20,
				Height:        180,
				WeightGoal:    "maintain",
				Sex:           "female",
				ActivityLevel: 5,
				Email:         "test_user@gmail.com",
			},
			want: errors.New("repository - user already exists in database"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := mockUserService.New(test.request)

			if !reflect.DeepEqual(err, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}
		})
	}
}

func TestGetAllUser(t *testing.T) {
	mockRepo := mockUserRepo{}
	mockUserService := api.NewUserService(&mockRepo)
	users := []api.User{
		api.User{
			ID:            1,
			Name:          "Rabbit",
			Age:           2,
			Height:        3,
			Sex:           "female",
			ActivityLevel: 2,
			WeightGoal:    "heavy",
			Email:         "light@email.com",
		},
		api.User{
			ID:            1,
			Name:          "Mole",
			Age:           2,
			Height:        3,
			Sex:           "female",
			ActivityLevel: 2,
			WeightGoal:    "heavy",
			Email:         "light2@email.com",
		},
	}

	tests := []struct {
		name       string
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
