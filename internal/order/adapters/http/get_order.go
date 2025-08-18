package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"r2-challenge/internal/order/services/query"
	"r2-challenge/pkg/observability"
)

type GetOrderHandler struct {
	service   query.GetByIDService
	validator *validator.Validate
	tracer    observability.Tracer
}

func NewGetOrderHandler(s query.GetByIDService, v *validator.Validate, t observability.Tracer) (GetOrderHandler, error) {
	return GetOrderHandler{service: s, validator: v, tracer: t}, nil
}

// Get Order by ID
// @Summary      Get order
// @Description  Get order by ID
// @Tags         Orders
// @Produce      json
// @Param        id   path     string  true  "Order ID"
// @Success      200  {object} map[string]any
// @Failure      400  {object} map[string]string "Bad Request"
// @Failure      404  {object} map[string]string "Not Found"
// @Router       /orders/{id} [get]
func (h GetOrderHandler) Handle(c echo.Context) error {
	ctx, span := h.tracer.StartSpan(c.Request().Context(), "OrderHTTP.GetByID")
	defer span.End()

	orderID := c.Param("id")
	if err := h.validator.Var(orderID, "required"); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	order, err := h.service.GetByID(ctx, orderID)
	if err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	// Authorization: only owner or admin can access
	role, _ := c.Get("role").(string)
	userID, _ := c.Get("userId").(string)
	if role != "admin" && order.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	return c.JSON(http.StatusOK, order)
}


