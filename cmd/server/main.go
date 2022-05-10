package main

import (
	"database/sql"
	"fmt"
	"os"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/app"
	"weight-tracker/pkg/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// if err := run(); err != nil
	// this syntax says:
	// 	set the value of err as the return value of run()
	//  and then do comparison, check if the value of err is not equal to nil
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup erorr: %s \\n", err)
		os.Exit(1)
	}
}

// func run will be responsible for setting up db connections, routers etc
func run() error {
	username := "chester"
	password := "baba_yetu"
	host := "localhost"
	databaseName := "weight_tracker"
	// I'm used to working with postgres, but feel free to use any db you like. You just have to change the driver
	// I'm not going to cover how to create a database here but create a database
	// and call it something along the lines of "weight trackerear
	// connectionString := fmt.Sprintf("password=baba_yetu dbname=%s sslmode=disable", databaseName)
	// connectionString := fmt.Sprintf("postgres://chester:baba_yetu@localhost/%s?sslmode=disable", databaseName)
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		username, password, host, databaseName,
	)

	// setup database connection
	db, err := setupDatabase(connectionString)

	if err != nil {
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(db)

	// create router dependecy
	router := gin.Default()
	router.Use(cors.Default())

	// create user service
	userService := api.NewUserService(storage)

	// create weight service
	weightService := api.NewWeightService(storage)

	// everything stays the same, so add this below
	// storage := repository.NewStorage(db)
	// run migrations

	// note that we are passing the connectionString again here. This is so
	// we can easily run migrations against another database, say a test version,
	// for our integration- and end-to-end tests
	err = storage.RunMigrations(connectionString)

	if err != nil {
		return err
	}

	server := app.NewServer(router, userService, weightService)

	// start the server
	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}

func setupDatabase(connString string) (*sql.DB, error) {
	// change "postgres" for whatever supported database you want to use
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	// ping the DB to ensure that it is connected
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
