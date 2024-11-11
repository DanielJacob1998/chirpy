package database

import (
    "time"
    "context"
)

func (q *Queries) InsertUser(ctx context.Context, email, hashedPassword string) (string, error) {
    var userID string
    
    // Current timestamp for created_at and updated_at
    now := time.Now()
    
    // Inserting user's email, hashed password, created_at, and updated_at
    err := q.db.QueryRowContext(ctx, `
        INSERT INTO users (email, hashed_password, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`, email, hashedPassword, now, now).Scan(&userID)
    
    if err != nil {
        return "", err
    }
    return userID, nil
}
