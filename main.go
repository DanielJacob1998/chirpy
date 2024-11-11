package main

import (
    "log"
    "net/http"
)

func main() {
    const filepathRoot = "."
    const port = "8080"

    // Use NewAPIConfig instead of manual creation
    apiCfg := NewAPIConfig()

    mux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    mux.Handle("/app/", fsHandler)

    // Update handler patterns with quotes
    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
    mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
    mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
    mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
    mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
    mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}
