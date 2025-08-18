package http

import (
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"

    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/services/query"
    "r2-challenge/pkg/observability"
)

type ListUsersHandler struct {
    service   query.ListService
    validator *validator.Validate
    tracer    observability.Tracer
}

func NewListUsersHandler(s query.ListService, v *validator.Validate, t observability.Tracer) (ListUsersHandler, error) {
    return ListUsersHandler{service: s, validator: v, tracer: t}, nil
}

// List Users
// @Summary      List users
// @Description  List users with optional filters
// @Tags         Users
// @Produce      json
// @Param        email   query    string  false  "Email"
// @Param        name    query    string  false  "Name contains"
// @Param        limit   query    int     false  "Limit"
// @Param        offset  query    int     false  "Offset"
// @Success      200     {array}  map[string]any
// @Failure      500     {object} map[string]string "Internal Server Error"
// @Router       /users [get]
func (h ListUsersHandler) Handle(c echo.Context) error {
    ctx, span := h.tracer.StartSpan(c.Request().Context(), "UserHTTP.List")
    defer span.End()

    f := repo.UserFilter{}
    if s := c.QueryParam("email"); s != "" { f.Email = s }
    if s := c.QueryParam("name"); s != "" { f.Name = s }
    if s := c.QueryParam("limit"); s != "" { if v, err := strconv.Atoi(s); err == nil { f.Limit = v } }
    if s := c.QueryParam("offset"); s != "" { if v, err := strconv.Atoi(s); err == nil { f.Offset = v } }

    list, err := h.service.List(ctx, f)
    if err != nil {
        span.RecordError(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, list)
}


