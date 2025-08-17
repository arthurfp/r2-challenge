package db

import (
    "context"
    "strings"
    "time"

    "gorm.io/gorm"

    "r2-challenge/internal/product/domain"
    appdb "r2-challenge/pkg/db"
    "r2-challenge/pkg/observability"
)

type dbProductRepository struct {
    db     *gorm.DB
    tracer observability.Tracer
}

func NewDBRepository(database *appdb.Database, t observability.Tracer) (ProductRepository, error) {
    return &dbProductRepository{db: database.DB, tracer: t}, nil
}

func (r *dbProductRepository) Save(ctx context.Context, p domain.Product) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Save")
    defer span.End()

    now := time.Now().UTC()
    p.CreatedAt = now
    p.UpdatedAt = now

    if err := r.db.WithContext(ctx).Table("products").Create(&p).Error; err != nil {
        span.RecordError(err)
        return domain.Product{}, err
    }
    return p, nil
}

func (r *dbProductRepository) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Update")
    defer span.End()

    p.UpdatedAt = time.Now().UTC()
    tx := r.db.WithContext(ctx).Table("products").Where("id = ?", p.ID).Updates(map[string]any{
        "name":        p.Name,
        "description": p.Description,
        "category":    p.Category,
        "price_cents": p.PriceCents,
        "inventory":   p.Inventory,
        "updated_at":  p.UpdatedAt,
    })
    if tx.Error != nil {
        span.RecordError(tx.Error)
        return domain.Product{}, tx.Error
    }
    if tx.RowsAffected == 0 {
        span.RecordError(gorm.ErrRecordNotFound)
        return domain.Product{}, gorm.ErrRecordNotFound
    }
    return r.GetByID(ctx, p.ID)
}

func (r *dbProductRepository) Delete(ctx context.Context, id string) error {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Delete")
    defer span.End()

    tx := r.db.WithContext(ctx).Table("products").Where("id = ?", id).Delete(nil)
    if tx.Error != nil {
        span.RecordError(tx.Error)
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        span.RecordError(gorm.ErrRecordNotFound)
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *dbProductRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.GetByID")
    defer span.End()

    var p domain.Product
    err := r.db.WithContext(ctx).Table("products").Where("id = ?", id).First(&p).Error
    if err != nil {
        span.RecordError(err)
        return domain.Product{}, err
    }
    return p, nil
}

func (r *dbProductRepository) List(ctx context.Context, f ProductFilter) ([]domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.List")
    defer span.End()

    var list []domain.Product
    q := r.db.WithContext(ctx).Table("products")

    if f.Category != "" {
        q = q.Where("category = ?", f.Category)
    }
    if f.Name != "" {
        q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(f.Name)+"%")
    }

    switch f.SortBy {
    case "price":
        q = q.Order("price_cents " + order(f.SortDesc))
    case "created_at":
        q = q.Order("created_at " + order(f.SortDesc))
    default:
        q = q.Order("name " + order(f.SortDesc))
    }

    if f.Limit > 0 {
        q = q.Limit(f.Limit)
    }
    if f.Offset > 0 {
        q = q.Offset(f.Offset)
    }

    if err := q.Find(&list).Error; err != nil {
        span.RecordError(err)
        return nil, err
    }
    return list, nil
}

func order(desc bool) string {
    if desc {
        return "desc"
    }
    return "asc"
}


