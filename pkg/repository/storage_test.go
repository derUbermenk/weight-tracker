package repository_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"weight-tracker/pkg/api"
)

var (
	db                *sql.DB
	connection_string string
)

const (
	username     = "chester"
	password     = "baba_yetu"
	host         = "localhost"
	databaseName = "weight_tracker_test"
)

func TestMain(m *testing.M) {
	formatConnectionString()
	connect_to_database(connection_string)
	set_up_database(connection_string)
	populate_database(connection_string)

	exitVal := m.Run() // run all test functions

	teardown_database(connection_string)

	os.Exit(exitVal) //

	// the above comments mean that I still need to reset the database after every test
	// but not necessarily set them up again.
}

func TestCreateUser(t *testing.T) {
	// defer DatabaseReset // this will only be called after the test is finished
	// define tests

	tests := []struct {
		name        string
		userRequest api.NewUserRequest
		want_uID    int
		want_error  error
	}{
		{
			name:        "Must return the user ID of the new user when successfully created",
			userRequest: api.NewUserRequest{},
			want_uID:    3,
			want_error:  nil,
		},
	}

	// run tests
	tx, err := db.Begin()
	// will be called after this function ends.
	// 	No return statement means function returns at end
	defer fmt.Println("called rollback") // delete this later
	defer tx.Rollback()

	if err != nil {
		t.Errorf("Error in initializing transaction: %v", err)
	}

	// make transaction
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uID, err := userRepo.CreateUser(test.userRequest)

			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if uID != test.want_uID {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, uID, test.want_uID)
			}
		})
	}
}
