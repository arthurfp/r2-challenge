package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/user/domain"
    "r2-challenge/internal/user/services/command"
    "r2-challenge/pkg/auth"
    "r2-challenge/pkg/observability"
)

type UpdateProfileHandler struct {
    service   command.UpdateProfileService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewUpdateProfileHandler(s command.UpdateProfileService, v *validator.Validate, t observability.Tracer) (UpdateProfileHandler, error) {
    return UpdateProfileHandler{service: s, validator: v, tracer: t}, nil
}

type updateProfileRequest struct {
    Name  string `json:"name" validate:"required,min=3"`
    Email string `json:"email" validate:"required,email"`
}

func (h UpdateProfileHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "UserHTTP.UpdateProfile")
    defer span.End()

    userID, _ := c.Get(auth.CtxUserID).(string)
    if err := h.validator.Var(userID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }

    var req updateProfileRequest
    if err := c.Bind(&req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }
	
    if err := h.validator.Struct(req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    user := domain.User{ID: userID, Name: req.Name, Email: req.Email}

    updated, err := h.service.UpdateProfile(ctx, user)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, updated)
}


