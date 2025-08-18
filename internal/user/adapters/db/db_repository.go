package db

import (
    "context"
    "strings"
    "time"

    "gorm.io/gorm"

    "r2-challenge/internal/user/domain"
    appdb "r2-challenge/pkg/db"
    "r2-challenge/pkg/observability"
    "github.com/google/uuid"
)

type dbUserRepository struct {
    db     *gorm.DB
    tracer observability.Tracer
}

func NewDBRepository(database *appdb.Database, t observability.Tracer) (UserRepository, error) {
    return &dbUserRepository{db: database.DB, tracer: t}, nil
}

func (r *dbUserRepository) Save(ctx context.Context, u domain.User) (domain.User, error) {
    ctx, span := r.tracer.StartSpan(ctx, "UserRepository.Save")
    defer span.End()

    now := time.Now().UTC()
    if u.ID == "" {
        u.ID = uuid.NewString()
    }
    u.CreatedAt = now
    u.UpdatedAt = now
    if err := r.db.WithContext(ctx).Table("users").Create(&u).Error; err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }
    return u, nil
}

func (r *dbUserRepository) Update(ctx context.Context, u domain.User) (domain.User, error) {
    ctx, span := r.tracer.StartSpan(ctx, "UserRepository.Update")
    defer span.End()

    u.UpdatedAt = time.Now().UTC()
    tx := r.db.WithContext(ctx).Table("users").Where("id = ?", u.ID).Updates(map[string]any{
        "email": u.Email,
        "name": u.Name,
        "role": u.Role,
        "password_hash": u.PasswordHash,
        "updated_at": u.UpdatedAt,
    })
    if tx.Error != nil { span.RecordError(tx.Error); return domain.User{}, tx.Error }
    if tx.RowsAffected == 0 { span.RecordError(gorm.ErrRecordNotFound); return domain.User{}, gorm.ErrRecordNotFound }
    return r.GetByID(ctx, u.ID)
}

func (r *dbUserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
    ctx, span := r.tracer.StartSpan(ctx, "UserRepository.GetByID")
    defer span.End()
    var u domain.User
    if err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).First(&u).Error; err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }
    return u, nil
}

func (r *dbUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
    ctx, span := r.tracer.StartSpan(ctx, "UserRepository.GetByEmail")
    defer span.End()
    var u domain.User
    if err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).First(&u).Error; err != nil {
        span.RecordError(err)
        return domain.User{}, err
    }
    return u, nil
}

func (r *dbUserRepository) List(ctx context.Context, f UserFilter) ([]domain.User, error) {
    ctx, span := r.tracer.StartSpan(ctx, "UserRepository.List")
    defer span.End()
    var list []domain.User
    q := r.db.WithContext(ctx).Table("users")
    if f.Email != "" { q = q.Where("email = ?", f.Email) }
    if f.Name != "" { q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(f.Name)+"%") }
    if f.Limit > 0 { q = q.Limit(f.Limit) }
    if f.Offset > 0 { q = q.Offset(f.Offset) }
    if err := q.Find(&list).Error; err != nil { span.RecordError(err); return nil, err }
    return list, nil
}


