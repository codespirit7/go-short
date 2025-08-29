package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codespirit7/url-shortner/controller"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()
	godotenv.Load()

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	c := cors.New(cors.Options{
    AllowedOrigins:   []string{"https://go-short.netlify.app"}, // better than "*"
    AllowCredentials: true,
    AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders:   []string{"*"}, // allow custom headers (Content-Type, Authorization, etc.)
    Debug: true, // turn on while debugging
})

	handler := c.Handler(mux)

	mux.HandleFunc("/short-url", controller.HandleShorten)
	mux.HandleFunc("/short/", controller.HandleRedirect)
	mux.HandleFunc("/health", controller.HandleHealth)

	fmt.Printf("URL Shortener is running on :%s", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), handler)

	if err != nil {
		log.Fatal(err)
	}

}
