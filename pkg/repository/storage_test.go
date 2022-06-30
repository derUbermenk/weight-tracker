package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/repository"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
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

// formats the connection string
func formatConnectionString() {
	connection_string = fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		username, password, host, databaseName,
	)
}

func connect_to_database(connection_string string) error {
	var err error
	db, err = sql.Open("postgres", connection_string)

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func set_up_database(connection_string string) error {
	if connection_string == "" {
		err := errors.New("repository: the connString was empty")
		return err
	}

	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, connection_string)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case errors.New("no change"):
		//
	default:
		return err
	}

	return nil
}

func populate_database(db *sql.DB) error {
	// use fixtures here

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("testdata/fixtures"),
	)

	if err != nil {
		return err
	}

	if err = fixtures.Load(); err != nil {
		return err
	}

	return nil
}

func teardown_database(connection_string string) error {
	// run down migrations instaed
	if connection_string == "" {
		err := errors.New("repository: the connString was empty")
		log.Printf("Error on teardown: %v", err)
		return err
	}

	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, connection_string)

	if err != nil {
		log.Printf("Error on teardown: %v", err)
		return err
	}

	err = m.Drop()

	switch err {
	case errors.New("no change"):
		//
	default:
		log.Printf("Error on teardown: %v", err)
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	var err error
	formatConnectionString()
	if err = connect_to_database(connection_string); err != nil {
		return
	}
	if err = set_up_database(connection_string); err != nil {
		return
	}
	if err = populate_database(db); err != nil {
		return
	}
	defer teardown_database(connection_string)

	exitVal := m.Run() // run all test functions
	defer os.Exit(exitVal)

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
	userRepo := repository.NewStorage(db)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// make transaction
			tx, err := db.Begin()
			if err != nil {
				t.Errorf("Error in initializing transaction: %v", err)
			}
			// will be called after this function ends.
			// 	No return statement means function returns at end
			defer fmt.Println("called rollback") // delete this later
			defer tx.Rollback()

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
