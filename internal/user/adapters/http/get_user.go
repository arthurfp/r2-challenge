package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/user/services/query"
    "r2-challenge/pkg/observability"
)

type GetUserHandler struct {
    service   query.GetByIDService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewGetUserHandler(s query.GetByIDService, v *validator.Validate, t observability.Tracer) (GetUserHandler, error) {
    return GetUserHandler{service: s, validator: v, tracer: t}, nil
}

func (h GetUserHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "UserHTTP.GetUser")
    defer span.End()

    id := c.Param("id")
    if err := h.validator.Var(id, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }

    u, err := h.service.GetByID(ctx, id)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
    }

    return c.JSON(http.StatusOK, u)
}


