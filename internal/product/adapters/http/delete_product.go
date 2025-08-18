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

// Delete Product
// @Summary      Delete product
// @Description  Delete a product by ID
// @Tags         Products
// @Produce      json
// @Param        id   path     string  true  "Product ID"
// @Success      204  {string} string  "No Content"
// @Failure      400  {object} map[string]string "Bad Request"
// @Failure      404  {object} map[string]string "Not Found"
// @Router       /v1/products/{id} [delete]
func (h DeleteHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.Delete")
    defer span.End()

    productID := c.Param("id")
    if err := h.validator.Var(productID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }

    if err := h.service.Delete(ctx, productID); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
    }

    return c.NoContent(http.StatusNoContent)
}


