package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

type Queries struct {
    db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
    return &Queries{db: db}
}

// Function to retrieve a single chirp by its ID
func (q *Queries) GetChirpByID(id string) (*Chirp, error) {
    chirp := &Chirp{}
    err := q.db.QueryRow("SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1", id).Scan(&chirp.ID, &chirp.CreatedAt, &chirp.UpdatedAt, &chirp.Body, &chirp.UserID)
    if err != nil {
        return nil, err
    }
    return chirp, nil
}
