package main

import (
    "encoding/json"
    "net/http"
    "time"
    "context"
    
    "github.com/DanielJacob1998/chirpy/internal/auth"
    "github.com/DanielJacob1998/chirpy/internal/database"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Password string `json:"password"`
        Email    string `json:"email"`
    }
    type response struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email     string    `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
        return
    }
    
    // Create the user using the CreateUser function
    userObj := database.CreateUserParams{
        Email:         params.Email,
        HashedPassword: hashedPassword,
    }
    user, err := cfg.CreateUser(r.Context(), userObj)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't insert user", err)
        return
    }
    
    // Construct the response using user details
    respondWithJSON(w, http.StatusCreated, response{
        ID:        user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email:     user.Email,
    })
}

func (cfg *apiConfig) CreateUser(ctx context.Context, params database.CreateUserParams) (User, error) {
    // Make sure InsertUser belongs to *database.Queries
    id, err := cfg.db.InsertUser(ctx, params.Email, params.HashedPassword)
    if err != nil {
        return User{}, err
    }
    
    user := User{
        ID:        uuid.MustParse(id), // Ensure the id is a valid UUID
        CreatedAt: time.Now(),         // Consider fetching actual creation time if possible
        UpdatedAt: time.Now(),
        Email:     params.Email,
    }
    return user, nil
}
