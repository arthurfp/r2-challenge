package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/product/services/command"
    "r2-challenge/pkg/observability"
)

type DeleteHandler struct {
    service   command.DeleteService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewDeleteHandler(s command.DeleteService, v *validator.Validate, t observability.Tracer) (DeleteHandler, error) {
    return DeleteHandler{service: s, validator: v, tracer: t}, nil
}

func (h DeleteHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.Delete")
    defer span.End()

    id := c.Param("id")
    if err := h.validator.Var(id, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }

    if err := h.service.Delete(ctx, id); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
    }

    return c.NoContent(http.StatusNoContent)
}


