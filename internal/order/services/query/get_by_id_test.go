package query

import (
    "context"
    "testing"

    gomock "github.com/golang/mock/gomock"
    orderdb "r2-challenge/internal/order/adapters/db"
    "r2-challenge/internal/order/domain"
    "r2-challenge/pkg/observability"
)

func TestOrderGetByID_Success(t *testing.T) {
    tracer, _ := observability.SetupTracer()
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)

    repo := orderdb.NewMockOrderRepository(ctrl)
    getSvc, _, _ := NewService(repo, tracer)

    repo.EXPECT().GetByID(gomock.Any(), "o1").Return(domain.Order{ID: "o1"}, nil)

    res, err := getSvc.GetByID(context.Background(), "o1")
    if err != nil { t.Fatalf("err: %v", err) }
    if res.ID != "o1" { t.Fatalf("id mismatch") }
}


