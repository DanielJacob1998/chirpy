package main

import (
    "database/sql" // Add this
    "encoding/json"
    "net/http"
    "strings"
    "time"
    "github.com/google/uuid"
)

type apiConfig struct {
    db       *database.Queries  // Change this line
    badWords map[string]struct{}
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    type createChirpRequest struct {
        Body   string `json:"body"`
        UserID string `json:"user_id"`
    }
    
    decoder := json.NewDecoder(r.Body)
    params := createChirpRequest{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // 2. Validate body (port your previous validation logic here)
    if len(params.Body) > 140 {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }
    
    // 3. Convert string UserID to UUID
    userID, err := uuid.Parse(params.UserID)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    // 4. Create database parameters
    chirpParams := database.CreateChirpParams{
        ID:        uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID:    userID,
        Body:      params.Body,
    }

    // 5. Insert into database
    chirp, err := cfg.DB.CreateChirp(r.Context(), chirpParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
        return
    }

    // 6. Respond with 201 and the chirp
    respondWithJSON(w, http.StatusCreated, chirp)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
    words := strings.Split(body, " ")
    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := badWords[loweredWord]; ok {
            words[i] = "****"
        }
    }
    return strings.Join(words, " ")
}
