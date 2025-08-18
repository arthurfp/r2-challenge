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

    "r2-challenge/internal/product/domain"
    "r2-challenge/pkg/observability"
    vsetup "r2-challenge/pkg/validator"
)

type fakeCreateService struct{ resp domain.Product; err error }
func (f fakeCreateService) Create(_ context.Context, p domain.Product) (domain.Product, error) { return f.resp, f.err }

func TestCreateProductHandler_ReturnsCreatedPayload(t *testing.T) {
    e := echo.New()
    v, _ := vsetup.Setup()
    tracer, _ := observability.SetupTracer()

    created := domain.Product{ID: "pid", Name: "N", Category: "C", PriceCents: 100}
    svc := fakeCreateService{resp: created}

    // build handler
    h, err := NewCreateHandler(svc, v, tracer)
    require.NoError(t, err)

    body := map[string]any{"name":"Name","description":"","category":"C","price_cents":100,"inventory":0}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/v1/products", bytes.NewReader(b))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err = h.Handle(c)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, rec.Code)

    var got domain.Product
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
    require.Equal(t, "pid", got.ID)
}


