package command

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	productdb "r2-challenge/internal/product/adapters/db"
	"r2-challenge/internal/product/domain"
	"r2-challenge/pkg/observability"
)

func TestCreateProductService_Success(t *testing.T) {
	tracer, _ := observability.SetupTracer()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := productdb.NewMockProductRepository(ctrl)

	service, err := NewCreateService(repo, tracer)
	if err != nil {
		t.Fatalf("unexpected error creating service: %v", err)
	}

	input := domain.Product{Name: "Item", Category: "cat", PriceCents: 1000, Inventory: 5}

	repo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p domain.Product) (domain.Product, error) {
		p.ID = "prod_1"
		return p, nil
	})

	result, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("service returned error: %v", err)
	}

	if result.ID == "" {
		t.Fatalf("expected ID to be set")
	}
}
