package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "sync/atomic"

    "github.com/gorilla/mux"
    "github.com/bootdotdev/learn-http-servers/internal/database"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

type apiConfig struct {
    fileserverHits atomic.Int32
    db             *database.Queries
    platform       string
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
    dbQueries := database.New(dbConn)

    apiCfg := apiConfig{
        fileserverHits: atomic.Int32{},
        db:             dbQueries,
        platform:       platform,
    }
    
    router := mux.NewRouter()

    // Static file server
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    router.PathPrefix("/app/").Handler(fsHandler)

    // API Routes
    router.HandleFunc("/api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            vars := mux.Vars(r)
            chirpID := vars["chirpID"]
            apiCfg.handlerChirpsRetrieve(w, r, chirpID)
        return
    }
    if r.Method == http.MethodPost {
        apiCfg.handlerChirpsCreate(w, r)
        return
    }
})

    router.HandleFunc("/api/healthz", handlerReadiness)
    router.HandleFunc("/api/users", apiCfg.handlerUsersCreate)
    router.HandleFunc("/admin/reset", apiCfg.handlerReset)
    router.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: router,
    }
    
    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}
