package main

import (
    "database/sql"
    "context"
    _ "github.com/lib/pq"
)

type Queries struct {
    db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
    return &Queries{db: db}
}

// Function to retrieve a single chirp by its ID
func (q *Queries) GetChirpByID(ctx context.Context, id string) (*Chirp, error) {
    var chirp Chirp

    // Note: we do not reopen a db connection here, use the one in q.db
    row := q.db.QueryRowContext(ctx, "SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1", id)

    err := row.Scan(&chirp.ID, &chirp.CreatedAt, &chirp.UpdatedAt, &chirp.Body, &chirp.UserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return &chirp, nil
}
