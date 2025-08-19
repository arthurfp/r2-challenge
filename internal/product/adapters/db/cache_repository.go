package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"r2-challenge/internal/product/domain"
	"r2-challenge/pkg/cache"
	appdb "r2-challenge/pkg/db"
	"r2-challenge/pkg/observability"
)

type cachedProductRepository struct {
	baseRepository ProductRepository
	cacheClient    *cache.Client
	cacheTTL       time.Duration
	tracer         observability.Tracer
}

func NewRepository(database *appdb.Database, t observability.Tracer, c *cache.Client) (ProductRepository, error) {
	// fallback to db-only if cache is nil
	baseRepository, err := NewDBRepository(database, t)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return baseRepository, nil
	}

	return &cachedProductRepository{baseRepository: baseRepository, cacheClient: c, cacheTTL: 30 * time.Second, tracer: t}, nil
}

func (r *cachedProductRepository) Save(ctx context.Context, product domain.Product) (domain.Product, error) {
	saved, err := r.baseRepository.Save(ctx, product)
	if err == nil {
		_ = r.cacheClient.Del(ctx, r.keyByID(saved.ID))
		_ = r.cacheClient.Del(ctx, r.keyListPrefix())
	}
	return saved, err
}

func (r *cachedProductRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	updated, err := r.baseRepository.Update(ctx, product)
	if err == nil {
		_ = r.cacheClient.Del(ctx, r.keyByID(updated.ID))
		_ = r.cacheClient.Del(ctx, r.keyListPrefix())
	}
	return updated, err
}

func (r *cachedProductRepository) Delete(ctx context.Context, productID string) error {
	err := r.baseRepository.Delete(ctx, productID)
	if err == nil {
		_ = r.cacheClient.Del(ctx, r.keyByID(productID))
		_ = r.cacheClient.Del(ctx, r.keyListPrefix())
	}
	return err
}

func (r *cachedProductRepository) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	if r.cacheClient == nil {
		return r.baseRepository.GetByID(ctx, productID)
	}

	if data, _ := r.cacheClient.Get(ctx, r.keyByID(productID)); data != nil {
		var cachedProduct domain.Product
		if err := json.Unmarshal(data, &cachedProduct); err == nil {
			return cachedProduct, nil
		}
	}

	product, err := r.baseRepository.GetByID(ctx, productID)
	if err != nil {
		return product, err
	}

	if encoded, err := json.Marshal(product); err == nil {
		_ = r.cacheClient.Set(ctx, r.keyByID(productID), encoded, r.cacheTTL)
	}

	return product, nil
}

func (r *cachedProductRepository) List(ctx context.Context, filter ProductFilter) ([]domain.Product, error) {
	if r.cacheClient == nil {
		return r.baseRepository.List(ctx, filter)
	}

	cacheKey := r.keyList(filter)
	if data, _ := r.cacheClient.Get(ctx, cacheKey); data != nil {
		var cachedList []domain.Product
		if err := json.Unmarshal(data, &cachedList); err == nil {
			return cachedList, nil
		}
	}

	list, err := r.baseRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if encoded, err := json.Marshal(list); err == nil {
		_ = r.cacheClient.Set(ctx, cacheKey, encoded, r.cacheTTL)
	}

	return list, nil
}

func (r *cachedProductRepository) keyByID(id string) string { return fmt.Sprintf("product:id:%s", id) }
func (r *cachedProductRepository) keyListPrefix() string    { return "product:list:" }
func (r *cachedProductRepository) keyList(f ProductFilter) string {
	return fmt.Sprintf("%sC=%s|N=%s|L=%d|O=%d|S=%s|D=%t", r.keyListPrefix(), f.Category, f.Name, f.Limit, f.Offset, f.SortBy, f.SortDesc)
}
