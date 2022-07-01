package repository

// update your imports to look like this:
import (
	"database/sql"
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"weight-tracker/pkg/api"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage interface {
	RunMigrations(connectionString string) error
	CreateUser(request api.NewUserRequest) (userID int, err error)
	CreateWeightEntry(request api.Weight) error
	DeleteUser(userID int) (deletedUserID int, err error)
	UpdateUser(request api.UpdateUserRequest) (api.User, error)
	GetUser(userID int) (api.User, error)
	GetUsers() ([]api.User, error)
	GetUserByEmail(userEmail string) (api.User, error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

// add this below NewStorage
func (s *storage) RunMigrations(connectionString string) error {
	if connectionString == "" {
		return errors.New("repository: the connString was empty")
	}
	// get base path
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")

	migrationsPath := filepath.Join("file://", basePath, "/pkg/repository/migrations/")

	m, err := migrate.New(migrationsPath, connectionString)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case errors.New("no change"):
		return nil
	}

	return nil
}

func (s *storage) CreateUser(request api.NewUserRequest) (userID int, err error) {
	newUserStatement := `
		INSERT INTO "user" (name, age, hashed_password, height, sex, activity_level, email, weight_goal)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
		`
	err = s.db.QueryRow(newUserStatement, request.Name, request.Age, request.HashedPassword, request.Height, request.Sex, request.ActivityLevel, request.Email, request.WeightGoal).Scan(&userID)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return
	}

	return
}

func (s *storage) DeleteUser(userID int) (deletedUserID int, err error) {
	deleteUserStatement := `
	DELETE FROM "user" 
	WHERE id=$1
	RETURNING id ;
	`

	err = s.db.QueryRow(deleteUserStatement, userID).Scan(&deletedUserID)

	if err != nil {
		log.Printf("storage error - this was the error: %v", err.Error())
		return
	}

	return
}

func (s *storage) UpdateUser(request api.UpdateUserRequest) (user api.User, err error) {
	updateUserStatement := `
		UPDATE "user" 
		SET name = $2, age = $3, height = $4,
		sex = $5, activity_level = $6, email = $7, 
		weight_goal = $8, updated_at = $9 WHERE id = $1
		RETURNING id, name, age, height, sex, activity_level, email,	
		weight_goal
		;`

	updateTime := time.Now()

	err = s.db.QueryRow(updateUserStatement,
		request.ID, request.Name, request.Age,
		request.Height, request.Sex, request.ActivityLevel,
		request.Email, request.WeightGoal, updateTime,
	).Scan(
		&user.ID, &user.Name, &user.Age,
		&user.Height, &user.Sex, &user.ActivityLevel,
		&user.Email, &user.WeightGoal,
	)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return
	}

	return
}

func (s *storage) CreateWeightEntry(request api.Weight) error {
	newWeightStatement := `
		INSERT INTO weight (weight, user_id, bmr, daily_caloric_intake)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`

	var ID int
	err := s.db.QueryRow(newWeightStatement, request.Weight, request.UserID, request.BMR, request.DailyCaloricIntake).Scan(&ID)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}

	return nil
}

func (s *storage) GetUsers() (users []api.User, err error) {
	getAllUsersStatement := `
		SELECT id, name, age, height, sex, activity_level, email, weight_goal,
		created_at, updated_at
		FROM "user";
	`
	// query users here
	rows, err := s.db.Query(getAllUsersStatement)

	if err != nil {
		return
	}

	// scan each user and add to users
	for rows.Next() {
		user := api.User{}
		if err = rows.Scan(
			&user.ID, &user.Name, &user.Age,
			&user.Height, &user.Sex, &user.ActivityLevel,
			&user.Email, &user.WeightGoal,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return
		}
		users = append(users, user)
	}

	rows.Close()
	return
}

func (s *storage) GetUser(userID int) (api.User, error) {
	getUserStatement := `
		SELECT id, name, age, height, sex, activity_level, email, weight_goal FROM "user"
		where id=$1;
		`

	var user api.User
	err := s.db.QueryRow(getUserStatement, userID).Scan(&user.ID, &user.Name, &user.Age, &user.Height, &user.Sex, &user.ActivityLevel, &user.Email, &user.WeightGoal)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.User{}, err
	}

	return user, nil
}

// queries for a user with given email. Returns
func (s *storage) GetUserByEmail(userEmail string) (user api.User, err error) {
	getUserByEmailStatement := `
		SELECT id, name, age, height, sex, activity_level, email, weight_goal FROM "user"
		where email=$1;
		`

	err = s.db.QueryRow(getUserByEmailStatement, userEmail).Scan(&user.ID, &user.Name, &user.Age, &user.Height, &user.Sex, &user.ActivityLevel, &user.Email, &user.WeightGoal)

	// no user with the given email was found in this case
	if errors.Is(err, sql.ErrNoRows) {
		return api.User{}, nil
	} else if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.User{}, err
	}

	// return the queried user if it does exist
	return user, nil
}
