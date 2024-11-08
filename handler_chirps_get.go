package main

import (
    "net/http"
    "github.com/gorilla/mux"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
    // Extract the chirpID from the request URL.
    vars := mux.Vars(r)
    chirpID := vars["chirpID"]

    // Fetch the chirp from the database.
    dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error fetching chirp", err)
        return
    }

    if dbChirp == nil { // Or however you check for "not found"
        respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
        return
    }

    // Prepare the response with the single chirp.
    chirp := Chirp{
        ID:        dbChirp.ID,
        CreatedAt: dbChirp.CreatedAt,
        UpdatedAt: dbChirp.UpdatedAt,
        UserID:    dbChirp.UserID,
        Body:      dbChirp.Body,
    }

    // Respond with JSON containing the chirp.
    respondWithJSON(w, http.StatusOK, chirp)
}
