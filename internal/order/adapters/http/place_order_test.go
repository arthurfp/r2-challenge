package http

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/require"

    "r2-challenge/internal/order/domain"
    "r2-challenge/pkg/auth"
    "r2-challenge/pkg/observability"
    vsetup "r2-challenge/pkg/validator"
)

// We generate mocks for command services separately if needed. For now, a small fake is enough.
type fakePlaceService struct{ resp domain.Order; err error }
func (f fakePlaceService) Place(_ context.Context, o domain.Order) (domain.Order, error) { return f.resp, f.err }

func TestPlaceOrderHandler_ReturnsCreatedOrderWithItems(t *testing.T) {
    e := echo.New()
    v, _ := vsetup.Setup()
    tracer, _ := observability.SetupTracer()

    placed := domain.Order{ID: "o1", UserID: "u1", Status: "created", Items: []domain.OrderItem{{ProductID: "p1", Quantity: 1, PriceCents: 100}}}
    svc := fakePlaceService{resp: placed}

    h, err := NewPlaceOrderHandler(svc, v, tracer)
    require.NoError(t, err)

    body := map[string]any{"items": []map[string]any{{"product_id": "11111111-1111-1111-1111-111111111111", "quantity": 1, "price_cents": 100}}}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/v1/orders", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set(auth.CtxUserID, "u1")

    err = h.Handle(c)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, rec.Code)

    var got domain.Order
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
    require.Equal(t, "o1", got.ID)
    require.Len(t, got.Items, 1)
    require.Equal(t, "p1", got.Items[0].ProductID)
}


