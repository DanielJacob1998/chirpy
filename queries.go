package db

import (
    "context"
    "database/sql"
    // other necessary imports
)

// Consider creating a struct to encapsulate your database methods
type Queries struct {
    db *sql.DB
}

// Function to retrieve a single chirp by its ID
func (q *Queries) GetChirpByID(ctx context.Context, id string) (*Chirp, error) {
    var chirp Chirp

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
