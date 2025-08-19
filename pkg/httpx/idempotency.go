package httpx

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "r2-challenge/pkg/cache"
)

type idempotencyRecord struct {
    Status int
    Body   []byte
    Hash   string
}

// IdempotencyMiddleware ensures POST operations are idempotent using a Redis cache.
// Expects header Idempotency-Key. Stores response for a short TTL.
func IdempotencyMiddleware(cch *cache.Client, ttl time.Duration) echo.MiddlewareFunc {
    if cch == nil { return func(next echo.HandlerFunc) echo.HandlerFunc { return next } }
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if c.Request().Method != http.MethodPost {
                return next(c)
            }

            idemKey := c.Request().Header.Get("Idempotency-Key")
            if idemKey == "" {
                return next(c)
            }

            // Read body and compute hash
            body, _ := io.ReadAll(c.Request().Body)
            _ = c.Request().Body.Close()
            c.Request().Body = io.NopCloser(bytes.NewReader(body))

            hashBytes := sha256.Sum256(body)
            payloadHash := hex.EncodeToString(hashBytes[:])

            cacheKey := "idem:" + idemKey
            if b, _ := cch.Get(c.Request().Context(), cacheKey); b != nil {
                var rec idempotencyRecord
                if json.Unmarshal(b, &rec) == nil {
                    if rec.Hash == payloadHash {
                        return c.Blob(rec.Status, echo.MIMEApplicationJSON, rec.Body)
                    }
                    return c.JSON(http.StatusConflict, map[string]string{"error": "idempotency key conflict"})
                }
            }

            // Capture response
            // Record response by wrapping writer
            orig := c.Response().Writer
            rr := &bodyCapture{ResponseWriter: orig}
            c.Response().Writer = rr
            if err := next(c); err != nil { return err }

            // Store
            rec := idempotencyRecord{ Status: c.Response().Status, Body: rr.buf, Hash: payloadHash }
            if enc, err := json.Marshal(rec); err == nil {
                _ = cch.Set(c.Request().Context(), cacheKey, enc, ttl)
            }
            return nil
        }
    }
}

// bodyCapture captures response body
type bodyCapture struct { http.ResponseWriter; buf []byte }
func (b *bodyCapture) Write(p []byte) (int, error) { b.buf = append(b.buf, p...); return b.ResponseWriter.Write(p) }


