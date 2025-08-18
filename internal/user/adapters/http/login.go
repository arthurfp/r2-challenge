package http

import (
    "net/http"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"

    "r2-challenge/cmd/envs"
    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/pkg/auth"
    "r2-challenge/pkg/observability"
)

type LoginHandler struct {
    repo      repo.UserRepository
    validator *validator.Validate
    tracer    observability.Tracer
    envs      envs.Envs
}

func NewLoginHandler(r repo.UserRepository, v *validator.Validate, t observability.Tracer, e envs.Envs) (LoginHandler, error) {
    return LoginHandler{repo: r, validator: v, tracer: t, envs: e}, nil
}

type loginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// Login
// @Summary      Login
// @Description  Authenticate and receive a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  loginRequest  true  "Login credentials"
// @Success      200  {object} map[string]string "access_token"
// @Failure      400  {object} map[string]string "Bad Request"
// @Failure      401  {object} map[string]string "Unauthorized"
// @Router       /v1/auth/login [post]
func (h LoginHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "UserHTTP.Login")
    defer span.End()

    var req loginRequest
    if err := c.Bind(&req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }
    if err := h.validator.Struct(req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    u, err := h.repo.GetByEmail(ctx, req.Email)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
    }

    ttl, _ := time.ParseDuration(h.envs.JWTExpire)
    tm := auth.NewTokenManager(h.envs.JWTSecret, h.envs.JWTIssuer, ttl)
    token, err := tm.Generate(u.ID, map[string]any{"email": u.Email, "role": u.Role})
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
    }

    return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}


