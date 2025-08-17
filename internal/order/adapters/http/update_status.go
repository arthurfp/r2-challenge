package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"r2-challenge/internal/order/services/command"
	"r2-challenge/pkg/observability"
)

type UpdateStatusHandler struct {
	service   command.UpdateStatusService
	validator *validator.Validate
	tracer    observability.Tracer
}

func NewUpdateStatusHandler(s command.UpdateStatusService, v *validator.Validate, t observability.Tracer) (UpdateStatusHandler, error) {
	return UpdateStatusHandler{service: s, validator: v, tracer: t}, nil
}

type updateStatusRequest struct { Status string `json:"status" validate:"required"` }

func (h UpdateStatusHandler) Handle(c echo.Context) error {
	ctx, span := h.tracer.StartSpan(c.Request().Context(), "OrderHTTP.UpdateStatus")
	defer span.End()

	orderID := c.Param("id")
	if err := h.validator.Var(orderID, "required"); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req updateStatusRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	if err := h.validator.Struct(req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	order, err := h.service.UpdateStatus(ctx, orderID, req.Status)
	if err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}


