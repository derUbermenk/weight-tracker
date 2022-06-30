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
	default:
		return err
	}

	return nil
}

func (dm *DatabaseManager) test_populateDB() error {
	if err := dm.fixtureLoader.Load(); err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	var err error
	databaseManager = NewDatabaseManager(username, password, host, databaseName)
	if err = databaseManager.main_connectDB(); err != nil {
		log.Printf("Error at main connect: %v", err)
		os.Exit(1)
	}

	if err = databaseManager.main_initializeFixtureLoader(); err != nil {
		log.Printf("Error at main fixture loader: %v", err)
		os.Exit(1)
	}

	exitValue := m.Run()
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

			uID, err := userRepo.CreateUser(test.userRequest)

			if err != test.want_error {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, err, test.want_error)
			}

			if uID != test.want_uID {
				t.Errorf("test: %v failed.\n\tgot: %v\n\twanted: %v", test.name, uID, test.want_uID)
			}

			if err = databaseManager.test_teardownDB(); err != nil {
				t.Errorf("Error at test setup: %v", err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {

}
