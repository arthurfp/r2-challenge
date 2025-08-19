package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"r2-challenge/internal/product/domain"
	"r2-challenge/internal/product/services/command"
	"r2-challenge/pkg/observability"
)

type CreateHandler struct {
	service   command.CreateService
	validator *validator.Validate
	tracer    observability.Tracer
}

func NewCreateHandler(s command.CreateService, v *validator.Validate, t observability.Tracer) (CreateHandler, error) {
	return CreateHandler{service: s, validator: v, tracer: t}, nil
}

type createProductRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	PriceCents  int64  `json:"price_cents" validate:"required,gte=0"`
	Inventory   int64  `json:"inventory" validate:"gte=0"`
}

// Create Product
// @Summary      Create product
// @Description  Create a new product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      createProductRequest  true  "Product input"
// @Success      201      {object}  domain.Product
// @Failure      400      {object}  map[string]string  "Bad Request"
// @Failure      500      {object}  map[string]string  "Internal Server Error"
// @Router       /products [post]
func (h CreateHandler) Handle(c echo.Context) error {
	ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.Create")
	defer span.End()

	var req createProductRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.validator.Struct(req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prod := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		PriceCents:  req.PriceCents,
		Inventory:   req.Inventory,
	}

	created, err := h.service.Create(ctx, prod)
	if err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, created)
}
