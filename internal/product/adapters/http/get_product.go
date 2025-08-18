package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/product/services/query"
    "r2-challenge/pkg/observability"
)

type GetHandler struct {
    service   query.GetByIDService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewGetHandler(s query.GetByIDService, v *validator.Validate, t observability.Tracer) (GetHandler, error) {
    return GetHandler{service: s, validator: v, tracer: t}, nil
}

// Get Product by ID
// @Summary      Get product
// @Description  Get a product by ID
// @Tags         Products
// @Produce      json
// @Param        id   path     string  true  "Product ID"
// @Success      200  {object} map[string]any
// @Failure      400  {object} map[string]string "Bad Request"
// @Failure      404  {object} map[string]string "Not Found"
// @Router       /v1/products/{id} [get]
func (h GetHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.GetByID")
    defer span.End()

    productID := c.Param("id")
    if err := h.validator.Var(productID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }

    product, err := h.service.GetByID(ctx, productID)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
    }

    return c.JSON(http.StatusOK, product)
}


