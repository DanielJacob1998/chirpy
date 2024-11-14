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
    jwtSecret      string
    polkaKey       string
}

func NewAPIConfig() *apiConfig {
    err := godotenv.Load(".env")
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

    // Debug prints
    dir, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting working directory: %v", err)
    }
    log.Println("Current working directory:", dir)
    log.Println("Environment variables loaded:")
    log.Println("DB_URL:", os.Getenv("DB_URL"))
    log.Println("PLATFORM:", os.Getenv("PLATFORM"))
    log.Println("JWT_SECRET:", len(os.Getenv("JWT_SECRET")), "bytes")

    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }
    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("PLATFORM must be set")
    }
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET must be set")
    }
    polkaKey := os.Getenv("POLKA_KEY")
    if polkaKey == "" {
        log.Fatal("POLKA_KEY must be set")
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
        jwtSecret:      jwtSecret,
        polkaKey:       polkaKey
    }
}
