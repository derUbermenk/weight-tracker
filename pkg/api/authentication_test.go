package api_test

import (
	"testing"
	"weight-tracker/pkg/api"
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
		return
	}

	user.Email = email
	user.Password = password

	return
}

var users = map[string]string{
	"existing_email@email.com": "correct_password1234",
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
		mockRepo := mockAuthRepo{users}
		authService := api.NewAuthService(&mockRepo)

		t.Run(test.name, func(t *testing.T) {
			validity, err := authService.ValidateCredentials(test.credentials)

			if err != test.want_err {
				t.Errorf("test %v failed. got: %v, wanted: %v", test.name, err, test.want_err)
			}

			if validity == test.want_validity {
				t.Errorf("test %v failed. got: %v, wanted: %v", test.name, validity, test.want_validity)
			}
		})
	}

}
