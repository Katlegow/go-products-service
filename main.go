package main

import "os"

func main() {
	app := App{}

	app.Initialise(
		os.Getenv("APP_DB_USER"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	app.Run(":8080")
}
