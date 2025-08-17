package command

import (
    "context"

    "golang.org/x/crypto/bcrypt"

    repo "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
    "r2-challenge/pkg/observability"
)

type RegisterService interface {
    Register(ctx context.Context, user domain.User, plainPassword string) (domain.User, error)
}

type registerService struct {
    repo   repo.UserRepository
    tracer observability.Tracer
}

func NewRegisterService(r repo.UserRepository, t observability.Tracer) (RegisterService, error) {
    return &registerService{repo: r, tracer: t}, nil
}

func (s *registerService) Register(ctx context.Context, user domain.User, plainPassword string) (domain.User, error) {
    ctx, span := s.tracer.StartSpan(ctx, "UserCommand.Register")
    defer span.End()

    hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }

    user.PasswordHash = string(hashed)

    saved, err := s.repo.Save(ctx, user)
    if err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }

    return saved, nil
}


