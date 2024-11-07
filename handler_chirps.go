package main

import (
    "encoding/json"
    "net/http"
    "strings"
)

func handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body   string `json:"body"`
        UserID string `json:"user_id"`
    }
    
    type chirpResponse struct {
        ID        string    `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
        UserID    string    `json:"user_id"`
    }
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
