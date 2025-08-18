package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/user/domain"
    "r2-challenge/internal/user/services/command"
    "r2-challenge/pkg/observability"
)

type RegisterHandler struct {
    service   command.RegisterService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewRegisterHandler(s command.RegisterService, v *validator.Validate, t observability.Tracer) (RegisterHandler, error) {
    return RegisterHandler{service: s, validator: v, tracer: t}, nil
}

type registerRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=3"`
    Password string `json:"password" validate:"required,min=6"`
}

// Register User
// @Summary      Register user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      registerRequest  true  "User input"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  map[string]string "Bad Request"
// @Failure      500   {object}  map[string]string "Internal Server Error"
// @Router       /auth/register [post]
func (h RegisterHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "UserHTTP.Register")
    defer span.End()

    var req registerRequest
    if err := c.Bind(&req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }
    if err := h.validator.Struct(req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    u := domain.User{Email: req.Email, Name: req.Name, Role: "user"}

    saved, err := h.service.Register(ctx, u, req.Password)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, saved)
}


