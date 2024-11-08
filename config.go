package main

import (
    "database/sql"
    "log"
    "os"
    "sync/atomic"
    
    "github.com/joho/godotenv"
    "github.com/DanielJacob1998/chirpy/internal/database"
    _ "github.com/lib/pq"
)

type apiConfig struct {
    fileserverHits atomic.Int32
    db             *database.Queries
    platform       string
}

func NewAPIConfig() *apiConfig {
    godotenv.Load()

    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }
    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("PLATFORM must be set")
    }

    dbConn, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Error opening database: %s", err)
    }
    dbQueries := database.New(dbConn)

    return &apiConfig{
        fileserverHits: atomic.Int32{},
        db:             dbQueries,
        platform:       platform,
    }
}
