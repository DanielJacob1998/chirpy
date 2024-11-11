package main

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/DanielJacob1998/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

    type parameters struct {
        Password          string `json:"password"`
        Email            string `json:"email"`
        ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
    }

    type response struct {
        ID        string    `json:"id"`
        Email     string    `json:"email"`
        Token     string    `json:"token"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    // Default to 1 hour if not specified
    expiresIn := time.Hour
    if params.ExpiresInSeconds != nil {
        // Convert seconds to duration
        expSeconds := time.Duration(*params.ExpiresInSeconds) * time.Second
        // Cap at 1 hour
        if expSeconds > time.Hour {
            expiresIn = time.Hour
        } else {
            expiresIn = expSeconds
        }
    }

    token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
        return
    }

    respondWithJSON(w, http.StatusOK, response{
        ID:        user.ID.String(),
        Email:     user.Email,
        Token:     token,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    })
}
