package command

import (
    "context"
    "testing"

    gomock "github.com/golang/mock/gomock"
    userdb "r2-challenge/internal/user/adapters/db"
    "r2-challenge/internal/user/domain"
    "r2-challenge/pkg/observability"
)

func TestRegister_Success(t *testing.T) {
    tracer, _ := observability.SetupTracer()
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)

    repo := userdb.NewMockUserRepository(ctrl)
    svc, _ := NewRegisterService(repo, tracer)

    repo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, u domain.User) (domain.User, error) {
        u.ID = "u1"
        return u, nil
    })

    user := domain.User{Email: "e@e.com", Name: "name", Role: "user"}
    res, err := svc.Register(context.Background(), user, "secret123")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if res.ID == "" { t.Fatalf("expected id") }
}


