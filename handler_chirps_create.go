package main

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/bootdotdev/learn-http-servers/internal/database"
    "github.com/google/uuid"
)

type Chirp struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    UserID    uuid.UUID `json:"user_id"`
    Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body   string    `json:"body"`
        UserID string `json:"user_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    userID, err := uuid.Parse(params.UserID)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user_id", err)
        return
    }
    
    cleanedBody, err := validateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid chirp", err)
        return
    }
    
    // Remove the duplicate CreateChirp call and keep only this one
    chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
        Body:   cleanedBody,  // Use the cleaned body here
        UserID: userID,
    })
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        return
    }

    respondWithJSON(w, http.StatusCreated, Chirp{
        ID:        chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body:      chirp.Body,
        UserID:    chirp.UserID,
    })
}

func validateChirp(body string) (string, error) {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return "", errors.New("Chirp is too long")
    }

    badWords := map[string]struct{}{
        "kerfuffle": {},
        "sharbert":  {},
        "fornax":    {},
    }
    cleaned := getCleanedBody(body, badWords)
    return cleaned, nil
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
    // Get chirps from database
    dbChirps, err := cfg.db.GetChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
        return
    }

    // Convert database chirps to API chirps
    chirps := []Chirp{}
    for _, dbChirp := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        })
    }

    respondWithJSON(w, http.StatusOK, chirps)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
    words := strings.Split(body, " ")
    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := badWords[loweredWord]; ok {
            words[i] = "****"
        }
    }
    cleaned := strings.Join(words, " ")
    return cleaned
}
