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

func TestGetByID_Success(t *testing.T) {
	tracer, _ := observability.SetupTracer()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := productdb.NewMockProductRepository(ctrl)
	getSvc, _, err := NewService(repo, tracer)
	if err != nil {
		t.Fatalf("new service error: %v", err)
	}

	repo.EXPECT().GetByID(gomock.Any(), "p1").Return(domain.Product{ID: "p1", Name: "N", Category: "C"}, nil)

	product, err := getSvc.GetByID(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != "p1" {
		t.Fatalf("expected id p1, got %s", product.ID)
	}
}

func TestGetByID_Error(t *testing.T) {
	tracer, _ := observability.SetupTracer()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := productdb.NewMockProductRepository(ctrl)
	getSvc, _, _ := NewService(repo, tracer)

	repoErr := errors.New("not found")
	repo.EXPECT().GetByID(gomock.Any(), "missing").Return(domain.Product{}, repoErr)

	_, err := getSvc.GetByID(context.Background(), "missing")
	if err == nil {
		t.Fatalf("expected error")
	}
}
