package database

import (
    "github.com/google/uuid"
    "database/sql"
    _ "github.com/lib/pq"
)

func (q *Queries) InsertUser(email string, hashedPassword string) (string, error) {
    // Generate a new UUID
    userID := uuid.New().String()

    // Execute the SQL insertion with the UUID
    err := q.db.QueryRow(
        "INSERT INTO users (id, email, hashed_password, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id",
        userID, email, hashedPassword,
    ).Scan(&userID)

    if err != nil {
        return "", err
    }

    return userID, nil
}
