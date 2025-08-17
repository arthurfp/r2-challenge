package query

import (
    "context"

    "r2-challenge/internal/user/domain"
)

type GetByEmailService interface { GetByEmail(ctx context.Context, email string) (domain.User, error) }

func (s *service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserQuery.GetByEmail")
    defer span.End()
    u, err := s.repo.GetByEmail(ctx, email)
    if err != nil { span.RecordError(err); return domain.User{}, err }
    return u, nil
}


