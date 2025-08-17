package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
    secret []byte
    issuer string
    ttl    time.Duration
}

func NewTokenManager(secret, issuer string, ttl time.Duration) TokenManager {
    return TokenManager{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

func (tm TokenManager) Generate(subject string, claims map[string]any) (string, error) {
    now := time.Now().UTC()
    std := jwt.MapClaims{
        "iss": tm.issuer,
        "sub": subject,
        "iat": now.Unix(),
        "exp": now.Add(tm.ttl).Unix(),
    }
    for k, v := range claims { std[k] = v }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, std)
    return token.SignedString(tm.secret)
}


