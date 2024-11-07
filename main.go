package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    "github.com/bootdotdev/learn-http-servers/internal/database"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type apiConfig struct {
    db *sql.DB  // Change this to *sql.DB
}

func main() {
    const filepathRoot = "."
    const port = "8080"

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

    apiCfg := apiConfig{
        db: dbConn,
    }

    mux := http.NewServeMux()
    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

    mux.HandleFunc("/api/healthz", handlerReadiness)
    mux.HandleFunc("/api/chirps", apiCfg.handlerChirpsCreate)  // note the apiCfg.
    mux.HandleFunc("/api/users", apiCfg.handlerUsersCreate)
    mux.HandleFunc("/admin/reset", apiCfg.handlerReset)
    mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)
    
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}
