package query

import (
    "context"
    "errors"
    "testing"

    gomock "github.com/golang/mock/gomock"
    productdb "r2-challenge/internal/product/adapters/db"
    "r2-challenge/internal/product/domain"
    "r2-challenge/pkg/observability"
)

func TestList_Success(t *testing.T) {
    tracer, _ := observability.SetupTracer()
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)

    repo := productdb.NewMockProductRepository(ctrl)
    _, listSvc, _ := NewService(repo, tracer)

    repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]domain.Product{{ID: "p1"}}, nil)

    products, err := listSvc.List(context.Background(), productdb.ProductFilter{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(products) != 1 {
        t.Fatalf("expected 1 product")
    }
}

func TestList_Error(t *testing.T) {
    tracer, _ := observability.SetupTracer()
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)

    repo := productdb.NewMockProductRepository(ctrl)
    _, listSvc, _ := NewService(repo, tracer)

    repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

    _, err := listSvc.List(context.Background(), productdb.ProductFilter{})
    if err == nil {
        t.Fatalf("expected error")
    }
}


