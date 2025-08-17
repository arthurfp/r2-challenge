package query

import (
    "context"

    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
)

type ListService interface { List(ctx context.Context, f repo.UserFilter) ([]domain.User, error) }

func (s *service) List(ctx context.Context, f repo.UserFilter) ([]domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserQuery.List")
    defer span.End()
    list, err := s.repo.List(ctx, f)
    if err != nil { span.RecordError(err); return nil, err }
    return list, nil
}


