package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/product/domain"
    "r2-challenge/internal/product/services/command"
    "r2-challenge/pkg/observability"
)

type UpdateHandler struct {
    service   command.UpdateService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewUpdateHandler(s command.UpdateService, v *validator.Validate, t observability.Tracer) (UpdateHandler, error) {
    return UpdateHandler{service: s, validator: v, tracer: t}, nil
}

type updateProductRequest struct {
    Name        string `json:"name" validate:"required,min=3"`
    Description string `json:"description"`
    Category    string `json:"category" validate:"required"`
    PriceCents  int64  `json:"price_cents" validate:"required,gte=0"`
    Inventory   int64  `json:"inventory" validate:"gte=0"`
}

// Update Product
// @Summary      Update product
// @Description  Update a product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path     string                 true  "Product ID"
// @Param        product  body     updateProductRequest   true  "Product input"
// @Success      200      {object} domain.Product
// @Failure      400      {object} map[string]string "Bad Request"
// @Failure      404      {object} map[string]string "Not Found"
// @Router       /products/{id} [put]
func (h UpdateHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.Update")
    defer span.End()

    productID := c.Param("id")
    if err := h.validator.Var(productID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }

    var req updateProductRequest
    if err := c.Bind(&req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }
    if err := h.validator.Struct(req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    product := domain.Product{ID: productID, Name: req.Name, Description: req.Description, Category: req.Category, PriceCents: req.PriceCents, Inventory: req.Inventory}
    updated, err := h.service.Update(ctx, product)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
    }

    return c.JSON(http.StatusOK, updated)
}


