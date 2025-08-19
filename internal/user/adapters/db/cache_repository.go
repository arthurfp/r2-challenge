package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"r2-challenge/internal/user/domain"
	"r2-challenge/pkg/cache"
	appdb "r2-challenge/pkg/db"
	"r2-challenge/pkg/observability"
)

type cachedUserRepository struct {
	baseRepository UserRepository
	cacheClient    *cache.Client
	cacheTTL       time.Duration
	tracer         observability.Tracer
}

func NewRepository(database *appdb.Database, t observability.Tracer, c *cache.Client) (UserRepository, error) {
	baseRepository, err := NewDBRepository(database, t)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return baseRepository, nil
	}
	return &cachedUserRepository{baseRepository: baseRepository, cacheClient: c, cacheTTL: 30 * time.Second, tracer: t}, nil
}

func (r *cachedUserRepository) Save(ctx context.Context, user domain.User) (domain.User, error) {
	saved, err := r.baseRepository.Save(ctx, user)
	if err == nil {
		_ = r.cacheClient.Del(ctx, r.keyByID(saved.ID))
		if saved.Email != "" {
			_ = r.cacheClient.Del(ctx, r.keyByEmail(saved.Email))
		}
		_ = r.cacheClient.Del(ctx, r.keyListPrefix())
	}
	return saved, err
}

func (r *cachedUserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	updated, err := r.baseRepository.Update(ctx, user)
	if err == nil {
		_ = r.cacheClient.Del(ctx, r.keyByID(updated.ID))
		if updated.Email != "" {
			_ = r.cacheClient.Del(ctx, r.keyByEmail(updated.Email))
		}
		_ = r.cacheClient.Del(ctx, r.keyListPrefix())
	}
	return updated, err
}

func (r *cachedUserRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	if r.cacheClient == nil {
		return r.baseRepository.GetByID(ctx, userID)
	}
	if data, _ := r.cacheClient.Get(ctx, r.keyByID(userID)); data != nil {
		var user domain.User
		if err := json.Unmarshal(data, &user); err == nil {
			return user, nil
		}
	}
	user, err := r.baseRepository.GetByID(ctx, userID)
	if err != nil {
		return user, err
	}
	if encoded, err := json.Marshal(user); err == nil {
		_ = r.cacheClient.Set(ctx, r.keyByID(userID), encoded, r.cacheTTL)
	}
	return user, nil
}

func (r *cachedUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if r.cacheClient == nil {
		return r.baseRepository.GetByEmail(ctx, email)
	}
	if data, _ := r.cacheClient.Get(ctx, r.keyByEmail(email)); data != nil {
		var user domain.User
		if err := json.Unmarshal(data, &user); err == nil {
			return user, nil
		}
	}
	user, err := r.baseRepository.GetByEmail(ctx, email)
	if err != nil {
		return user, err
	}
	if encoded, err := json.Marshal(user); err == nil {
		_ = r.cacheClient.Set(ctx, r.keyByEmail(email), encoded, r.cacheTTL)
	}
	return user, nil
}

func (r *cachedUserRepository) List(ctx context.Context, filter UserFilter) ([]domain.User, error) {
	if r.cacheClient == nil {
		return r.baseRepository.List(ctx, filter)
	}
	cacheKey := r.keyList(filter)
	if data, _ := r.cacheClient.Get(ctx, cacheKey); data != nil {
		var list []domain.User
		if err := json.Unmarshal(data, &list); err == nil {
			return list, nil
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

func (r *cachedUserRepository) keyByID(id string) string { return fmt.Sprintf("user:id:%s", id) }
func (r *cachedUserRepository) keyByEmail(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}
func (r *cachedUserRepository) keyListPrefix() string { return "user:list:" }
func (r *cachedUserRepository) keyList(f UserFilter) string {
	return fmt.Sprintf("%sE=%s|N=%s|L=%d|O=%d", r.keyListPrefix(), f.Email, f.Name, f.Limit, f.Offset)
}
