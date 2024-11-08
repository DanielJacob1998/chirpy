package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    const filepathRoot = "."
    const port = "8080"
    
    // Initialize the API configuration
    apiCfg := NewAPIConfig()

    // Setup routes using the helper function
    router := setupRoutes(apiCfg, filepathRoot)

    // Create server
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: router,
    }

    // Start server
    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}

// SetupRoutes is below main for organization
func setupRoutes(apiCfg apiConfig, filepathRoot string) *mux.Router {
    router := mux.NewRouter()

    // File server
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    router.PathPrefix("/app/").Handler(fsHandler)

    // API Routes
    router.HandleFunc("/api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        chirpID := vars["chirpID"]

        switch r.Method {
        case http.MethodGet:
            apiCfg.handlerChirpsRetrieve(w, r, chirpID)
        case http.MethodPost:
            apiCfg.handlerChirpsCreate(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    router.HandleFunc("/api/healthz", apiCfg.handlerReadiness)
    router.HandleFunc("/api/users", apiCfg.handlerUsersCreate)
    router.HandleFunc("/admin/reset", apiCfg.handlerReset)
    router.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)

    return router
}
