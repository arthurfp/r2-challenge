package http

import (
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    repo "r2-challenge/internal/product/adapters/db"
    "r2-challenge/internal/product/services/query"
    "r2-challenge/pkg/observability"
)

type ListHandler struct {
    service   query.ListService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewListHandler(s query.ListService, v *validator.Validate, t observability.Tracer) (ListHandler, error) {
    return ListHandler{service: s, validator: v, tracer: t}, nil
}

func (h ListHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "ProductHTTP.List")
    defer span.End()

    f := repo.ProductFilter{
        Category: c.QueryParam("category"),
        Name:     c.QueryParam("name"),
        SortBy:   c.QueryParam("sort"),
        SortDesc: c.QueryParam("order") == "desc",
    }
    if s := c.QueryParam("limit"); s != "" {
        if v, err := strconv.Atoi(s); err == nil { f.Limit = v }
    }
    if s := c.QueryParam("offset"); s != "" {
        if v, err := strconv.Atoi(s); err == nil { f.Offset = v }
    }

    list, err := h.service.List(ctx, f)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, list)
}


