package query

import (
    "context"

    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
    "r2-challenge/pkg/observability"
)

type GetByEmailService interface { GetByEmail(ctx context.Context, email string) (domain.User, error) }

func NewGetByEmailService(r repo.UserRepository, t observability.Tracer) (GetByEmailService, error) {
    return &service{repo: r, tracer: t}, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserQuery.GetByEmail")
    defer span.End()
    u, err := s.repo.GetByEmail(ctx, email)
    if err != nil { span.RecordError(err); return domain.User{}, err }
    return u, nil
}


