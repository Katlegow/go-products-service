package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Defines our application configurations
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Creates db connection and wires up corresponding routers
func (app *App) Initialise(user, password, dbname string) {

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	app.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()

}

//Runs the application
func (app *App) Run(addr string) {

}
