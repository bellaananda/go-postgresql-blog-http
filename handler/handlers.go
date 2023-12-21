package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bellaananda/go-postgresql-blog-http.git/database"
)

func FirstHandler(w http.ResponseWriter, r *http.Request) {
	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, your API is working!"))
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize the database and repositories
	db, errDb := database.RunDatabase()
	if errDb != nil {
		msg := fmt.Sprintf(`{"error": "%s"}`, errDb.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Successfully connected!"))

	// close db connection
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err.Error()) // proper error handling
		}
		sqlDB.Close()
	}()
}

func MigrateHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize the database and repositories
	db, errDb := database.RunDatabase()
	if errDb != nil {
		msg := fmt.Sprintf(`{"error": "%s"}`, errDb.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	// Migrate the database
	dbIn := database.NewPostgreSQLGORMRepository(db)
	if err := dbIn.Migrate(context.Background()); err != nil {
		msg := fmt.Sprintf(`{"error": "%s"}`, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully migrated!"))

	// close db connection
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err.Error()) // proper error handling
		}
		sqlDB.Close()
	}()
}
