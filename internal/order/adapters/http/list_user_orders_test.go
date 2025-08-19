package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	odb "r2-challenge/internal/order/adapters/db"
	"r2-challenge/internal/order/domain"
	oq "r2-challenge/internal/order/services/query"
	"r2-challenge/pkg/auth"
	"r2-challenge/pkg/observability"
	vsetup "r2-challenge/pkg/validator"
)

func TestListUserOrdersHandler_ReturnsItems(t *testing.T) {
	e := echo.New()
	v, _ := vsetup.Setup()
	tracer, _ := observability.SetupTracer()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockSvc := oq.NewMockListByUserService(ctrl)

	handler, err := NewListUserOrdersHandler(mockSvc, v, tracer)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/any/orders?limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Inject authenticated user id into context
	c.Set(auth.CtxUserID, "user-1")

	expected := []domain.Order{{
		ID:         "o1",
		UserID:     "user-1",
		Status:     "created",
		TotalCents: 100,
		Items: []domain.OrderItem{{
			ID:         "i1",
			OrderID:    "o1",
			ProductID:  "p1",
			Quantity:   1,
			PriceCents: 100,
		}},
	}}

	mockSvc.EXPECT().ListByUser(gomock.Any(), "user-1", odb.OrderFilter{Limit: 10}).Return(expected, nil)

	err = handler.Handle(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var got []domain.Order
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Len(t, got[0].Items, 1)
	require.Equal(t, "p1", got[0].Items[0].ProductID)
}
