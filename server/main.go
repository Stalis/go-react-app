package main

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading \".env\" file " + err.Error())
	}
}

func main() {
	frontendPath := "../frontend/build"

	router := mux.NewRouter()

	//fileServer := http.FileServer(http.Dir(frontendPath))
	router.PathPrefix("/js")
}
