package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yonasadiel/helios"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading env")
	}

	helios.App.Initialize()

	defer helios.App.CloseDB()

	helios.App.Migrate()

	r := CreateRouter()
	fmt.Println("Starting server on port 8100...")
	log.Fatal(http.ListenAndServe(":8100", r))
}
