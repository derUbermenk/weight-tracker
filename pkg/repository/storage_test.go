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

const (
	username     = "chester"
	password     = "baba_yetu"
	host         = "localhost"
	databaseName = "weight_tracker_test"
)

var databaseManager *DatabaseManager

type DatabaseManager struct {
	db                *sql.DB
	fixtureLoader     *testfixtures.Loader
	connection_string string
}

func NewDatabaseManager(username, password, host, databaseName string) *DatabaseManager {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		username, password, host, databaseName,
	)

	return &DatabaseManager{
		connection_string: connectionString,
	}
}

func (dm *DatabaseManager) main_connectDB() error {
	var err error
	dm.db, err = sql.Open("postgres", dm.connection_string)

	if err != nil {
		return err
	}

	err = dm.db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) main_initializeFixtureLoader() error {
	var err error

	if dm.db == nil {
		err = errors.New("No database connected")
		return err
	}

	dm.fixtureLoader, err = testfixtures.New(
		testfixtures.Database(dm.db),       // You database connection
		testfixtures.Dialect("postgres"),   // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("fixtures"), // The directory containing the YAML files
	)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) test_setupDB() error {
	if dm.connection_string == "" {
		return errors.New("repository: the connString was empty")
	}
	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, dm.connection_string)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case migrate.ErrNoChange:
		return nil
	default:
		return err
	}

	return nil
}

func (dm *DatabaseManager) test_teardownDB() error {
	if dm.connection_string == "" {
		return errors.New("repository: the connString was empty")
	}
	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, dm.connection_string)

	if err != nil {
		return err
	}

	err = m.Drop()

	switch err {
	case migrate.ErrNoChange:
		return nil
	case nil:
		return nil
	default:
		return err
	}
}

func (dm *DatabaseManager) test_populateDB() error {
	if err := dm.fixtureLoader.Load(); err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) testSetUpGroup(t *testing.T) {
	var err error

	if err = databaseManager.test_teardownDB(); err != nil {
		t.Errorf("Error at test setup: %v", err)
	}

	if err = databaseManager.test_setupDB(); err != nil {
		t.Errorf("Error at test setup: %v", err)
	}

	if err = databaseManager.test_populateDB(); err != nil {
		t.Errorf("Error at test setup: %v", err)
	}
}

func (dm *DatabaseManager) testTearDownGroup(t *testing.T) {
	var err error
	if err = databaseManager.test_teardownDB(); err != nil {
		t.Errorf("Error at test setup: %v", err)
	}
}

func (dm *DatabaseManager) mainSetUpGroup() {
	// the success of the following tests and their reliance on fixtures
	// relies on the ordering to this main setup group.
	//
	// main_connectDB connects to the database and initilizes the db field of the database manager
	// test_setupDB relies on the db initialized by the prior function to setup the tables
	// fixtureloader initilizes a fixture loader it does not populate the db yet.
	//
	// apparently the fixture loader relies on the fact that the db has already tables in store
	// to ensure proper functionality. earlier I did not run test_setupDB prior and had failing tests
	// due to the fixture loader not knowing proper ordering

	var err error
	if err = dm.main_connectDB(); err != nil {
		log.Printf("Error at main setup:\n\tConnecting to db error: %v", err)
		os.Exit(1)
	}

	if err = dm.test_setupDB(); err != nil {
		log.Printf("Error at main setup:\n\tSetting up db error: %v", err)
		os.Exit(1)
	}

	if err = dm.main_initializeFixtureLoader(); err != nil {
		log.Printf("Error at main setup:\n\tInitializing fixture loader error: %v", err)
		os.Exit(1)
	}
}

func (dm *DatabaseManager) mainTearDownGroup() {
	var err error
	if err = dm.test_teardownDB(); err != nil {
		log.Printf("Error at main teardown:\n\ttear down error: %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	databaseManager = NewDatabaseManager(username, password, host, databaseName)
	databaseManager.mainSetUpGroup()

	exitValue := m.Run()

	databaseManager.mainTearDownGroup()
	os.Exit(exitValue)
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
			userRequest: api.NewUserRequest{Name: "rabbit", Email: "rabbit@email.com", HashedPassword: "hashedPass"},
			want_uID:    10001,
			want_error:  nil,
		},
	}

	// run tests
	userRepo := repository.NewStorage(databaseManager.db)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			databaseManager.testSetUpGroup(t)

			uID, err := userRepo.CreateUser(test.userRequest)

			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if uID != test.want_uID {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, uID, test.want_uID)
			}

			databaseManager.testTearDownGroup(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name                      string
		delete_id                 int
		want_deleted_id           int
		want_error                error
		want_remaining_user_count int
	}{
		{
			name:       "Must not return an error",
			delete_id:  1,
			want_error: nil,
		},
		{
			name:            "Must return the correct id of the users",
			delete_id:       1,
			want_deleted_id: 1,
		},
		{
			name:                      "Must leave the correct amount of remaining users",
			delete_id:                 1,
			want_remaining_user_count: 1,
		},
	}

	userRepo := repository.NewStorage(databaseManager.db)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			databaseManager.testSetUpGroup(t)

			switch test.name {

			case tests[0].name:
				_, err := userRepo.DeleteUser(test.delete_id)

				if err != test.want_error {
					t.Errorf("test: %v failed. \n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
				}

			case tests[1].name:
				deleted_id, _ := userRepo.DeleteUser(test.delete_id)

				if deleted_id != test.delete_id {
					t.Errorf("test: %v failed. \n\tgot: %v\n\twanted: %v", test.name, deleted_id, test.want_deleted_id)
				}

			case tests[2].name:
				userRepo.DeleteUser(test.delete_id)

				var remaining_user_count int
				count_user_statement := `
					SELECT COUNT(*) FROM "user";
				`
				err := databaseManager.db.QueryRow(count_user_statement).Scan(&remaining_user_count)

				if err != nil {
					t.Errorf("test: %v failed. \n\tGot error: %v", test.name, err)
				}

				if remaining_user_count != test.want_remaining_user_count {
					t.Errorf("test: %v failed. \n\tgot: %v\n\twanted: %v", test.name, remaining_user_count, test.want_remaining_user_count)
				}
			}

			databaseManager.testTearDownGroup(t)
		})
	}
}
