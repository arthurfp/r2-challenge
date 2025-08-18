package query

import (
    "context"

    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
    "r2-challenge/pkg/observability"
)

type GetByIDService interface { GetByID(ctx context.Context, id string) (domain.User, error) }

type service struct {
    repo   repo.UserRepository
    tracer observability.Tracer
}

func NewGetByIDService(r repo.UserRepository, t observability.Tracer) (GetByIDService, error) {
    return &service{repo: r, tracer: t}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserQuery.GetByID")
    defer span.End()
    u, err := s.repo.GetByID(ctx, id)
    if err != nil { span.RecordError(err); return domain.User{}, err }
    return u, nil
}


