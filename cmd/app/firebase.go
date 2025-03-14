package app

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	// Determine credentials file path based on environment
	credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credPath == "" {
		credPath = "config/service-account-key.json" // Default local path
	}

	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatal("Error initializing Firebase: ", err)
	}

	FirebaseApp = app
	log.Println("Firebase initialized with credentials file:", credPath)
}
