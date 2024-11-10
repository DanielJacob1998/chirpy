package database

import (
    _ "github.com/lib/pq"
    "context"
    "github.com/DanielJacob1998/chirpy/internal/auth"
)

func (q *Queries) InsertUser(ctx context.Context, email, hashedPassword string) (string, error) {
    var userID string
    
    // Assume your users table has columns email, hashed_password, and id
    err := q.db.QueryRowContext(ctx, `
        INSERT INTO users (email, hashed_password)
        VALUES ($1, $2)
        RETURNING id`, email, hashedPassword).Scan(&userID)
    
    if err != nil {
        return "", err
    }
    return userID, nil
}
