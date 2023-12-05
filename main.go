package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Sysleec/Blog-Aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	//This is my .env file
	//you need to change it to tokens.env
	err := godotenv.Load("secret.env")

	if err != nil {
		log.Fatalf("You don't have this env file. %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	dbURL := os.Getenv("DB_URL")
	if port == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to DB:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.handlerGetUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
