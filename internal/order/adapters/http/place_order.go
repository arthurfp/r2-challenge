package http

import (
    "net/http"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    "r2-challenge/internal/order/domain"
    "r2-challenge/internal/order/services/command"
    "r2-challenge/pkg/auth"
    "r2-challenge/pkg/observability"
)

type PlaceOrderHandler struct {
    service   command.PlaceOrderService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewPlaceOrderHandler(s command.PlaceOrderService, v *validator.Validate, t observability.Tracer) (PlaceOrderHandler, error) {
    return PlaceOrderHandler{service: s, validator: v, tracer: t}, nil
}

type placeOrderRequest struct {
    Items []struct {
        ProductID  string `json:"product_id" validate:"required"`
        Quantity   int64  `json:"quantity" validate:"required,gt=0"`
        PriceCents int64  `json:"price_cents" validate:"required,gte=0"`
    } `json:"items" validate:"required,dive"`
}

func (h PlaceOrderHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "OrderHTTP.Place")
    defer span.End()

    var req placeOrderRequest
    if err := c.Bind(&req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }

    if err := h.validator.Struct(req); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    userID, _ := c.Get(auth.CtxUserID).(string)
    if err := h.validator.Var(userID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }

    items := make([]domain.OrderItem, 0, len(req.Items))
    var total int64
    for _, it := range req.Items {
        items = append(items, domain.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity, PriceCents: it.PriceCents})
        total += it.PriceCents * it.Quantity
    }

    ord := domain.Order{UserID: userID, Items: items, TotalCents: total, Status: "created"}
    saved, err := h.service.Place(ctx, ord)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
	return c.JSON(http.StatusCreated, saved)
}


