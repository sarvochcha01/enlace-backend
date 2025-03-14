package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/sarvochcha01/enlace-backend/internal/routes"
)

type App struct {
	router     chi.Router
	db         *sql.DB
	authClient *auth.Client
}

func (a *App) Initialise() {

	InitFirebase()
	var err error
	a.authClient, err = FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatal("Failed to initialize Firebase Auth: ", err)
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=123456 dbname=enlace sslmode=disable host=localhost port=5431"
	}

	a.db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Failed to open db: ", err)
		return
	}

	if err := a.db.Ping(); err != nil {
		log.Fatal("Error connecting to db: ", err)
		return
	}

	a.router = chi.NewRouter()

	a.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://enlace-frontend.vercel.app"}, // Frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                     // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},                               // Allowed headers
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight for 5 minutes
	}).Handler)

	routes.SetupRoutes(a.router, a.db, a.authClient)
}

func (a *App) Run() {
	log.Println("Server running on port 3000")
	if err := http.ListenAndServe(":3000", a.router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
