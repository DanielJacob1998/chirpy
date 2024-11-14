package auth

import (
    "fmt"
    "net/http"
    "strings"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "time"
    "crypto/rand"
    "encoding/hex"
)

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
    claims := jwt.RegisteredClaims{
        Issuer:    "chirpy",
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        Subject:   userID.String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := &jwt.RegisteredClaims{}
    
    // First, we parse the token
    token, err := jwt.ParseWithClaims(
        tokenString,
        claims,
        func(token *jwt.Token) (interface{}, error) {
            return []byte(tokenSecret), nil
        },
    )

    // Check for parsing errors
    if err != nil {
        return uuid.Nil, err
    }

    // Check if token is valid
    if !token.Valid {
        return uuid.Nil, fmt.Errorf("invalid token")
    }

    // Get the user ID from the Subject claim
    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
        return uuid.Nil, err
    }

    return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("Authorization header not found")
    }
    
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return "", fmt.Errorf("Authorization header must start with Bearer")
    }
    
    return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func MakeRefreshToken() (string, error) {
    slice := make([]byte, 32)
    _, err := rand.Read(slice)  // We pass our slice to rand.Read
    if err != nil {
        return "", err
    }
    hexxed := hex.EncodeToString(slice)
    return hexxed, nil
}

// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeaderIncluded
    }
    splitAuth := strings.Split(authHeader, " ")
    if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
        return "", errors.New("malformed authorization header")
    }

    return splitAuth[1], nil
}
