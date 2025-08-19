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

func NewListUsersHandler(service query.ListService, v *validator.Validate, tracer observability.Tracer) (ListUsersHandler, error) {
	return ListUsersHandler{service: service, validator: v, tracer: tracer}, nil
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

	filter := repo.UserFilter{}

	if emailParam := c.QueryParam("email"); emailParam != "" {
		filter.Email = emailParam
	}

	if nameParam := c.QueryParam("name"); nameParam != "" {
		filter.Name = nameParam
	}

	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			filter.Limit = limit
		}
	}

	if offsetParam := c.QueryParam("offset"); offsetParam != "" {
		if offset, err := strconv.Atoi(offsetParam); err == nil {
			filter.Offset = offset
		}
	}

	list, err := h.service.List(ctx, filter)
	if err != nil {
		span.RecordError(err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}
