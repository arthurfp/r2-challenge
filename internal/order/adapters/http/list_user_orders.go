package http

import (
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    repo "r2-challenge/internal/order/adapters/db"
    "r2-challenge/internal/order/services/query"
    "r2-challenge/pkg/auth"
    "r2-challenge/pkg/observability"
)

type ListUserOrdersHandler struct {
    service   query.ListByUserService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewListUserOrdersHandler(s query.ListByUserService, v *validator.Validate, t observability.Tracer) (ListUserOrdersHandler, error) {
    return ListUserOrdersHandler{service: s, validator: v, tracer: t}, nil
}

// List My Orders
// @Summary      List my orders
// @Description  List orders for the authenticated user
// @Tags         Orders
// @Produce      json
// @Param        id      path     string  true  "User ID"
// @Param        limit   query    int  false  "Limit"
// @Param        offset  query    int  false  "Offset"
// @Success      200     {array}  map[string]any
// @Failure      401     {object} map[string]string "Unauthorized"
// @Failure      500     {object} map[string]string "Internal Server Error"
// @Router       /users/{id}/orders [get]
func (h ListUserOrdersHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "OrderHTTP.ListByUser")
    defer span.End()

    userID, _ := c.Get(auth.CtxUserID).(string)
    if err := h.validator.Var(userID, "required"); err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }

    var filter repo.OrderFilter
    if s := c.QueryParam("limit"); s != "" {
        if v, err := strconv.Atoi(s); err == nil {
            filter.Limit = v
        }
    }
    if s := c.QueryParam("offset"); s != "" {
        if v, err := strconv.Atoi(s); err == nil {
            filter.Offset = v
        }
    }

    list, err := h.service.ListByUser(ctx, userID, filter)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
	return c.JSON(http.StatusOK, list)
}


