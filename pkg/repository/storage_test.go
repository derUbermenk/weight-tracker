package repository_test

import (
	"database/sql"
	"os"
	"testing"
	"weight-tracker/pkg/repository"
)

var userRepo repository.Storage
var db *sql.DB

// sets up the testing database
func setup_db() {
	userRepo = repository.NewStorage(db)
	userRepo.RunMigrations("")

	// run fixtures
}

func TestMain(m *testing.M) {
	setup_db()
	exitVal := m.Run() // run all test functions
	os.Exit(exitVal)   //

	// the above comments mean that I still need to reset the database after every test
	// but not necessarily set them up again.
}

func TestCreateUser(t *testing.T) {
	// defer DatabaseReset // this will only be called after the test is finished
	// define tests
	// run tests
	tx, err := db.Begin()

	if err != nil {
		t.Errorf("Error in initializing transaction: %v", err)
	}
	// will be called after this function ends.
	// 	No return statement means function returns at end
	defer tx.Rollback()

	// make transaction
}
