package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/yonasadiel/charon/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading env")
	}

	app.Charon.Initialize()

	defer app.Charon.CloseDB()

	app.Charon.Migrate()

	r := CreateRouter()
	fmt.Println("Starting server on port 8100...")
	log.Fatal(http.ListenAndServe(":8100", r))
}
