package command

import (
    "context"

    userdb "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
    "r2-challenge/pkg/observability"
)

type UpdateProfileService interface { UpdateProfile(ctx context.Context, user domain.User) (domain.User, error) }

type updateProfileService struct {
    repo   userdb.UserRepository
    tracer observability.Tracer
}

func NewUpdateProfileService(r userdb.UserRepository, t observability.Tracer) (UpdateProfileService, error) {
    return &updateProfileService{repo: r, tracer: t}, nil
}

func (s *updateProfileService) UpdateProfile(ctx context.Context, user domain.User) (domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserCommand.UpdateProfile")
    defer span.End()

    updated, err := s.repo.Update(ctx, user)
    if err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }

    return updated, nil
}


