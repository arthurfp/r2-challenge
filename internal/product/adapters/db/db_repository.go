package db

import (
    "context"
    "strings"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    "r2-challenge/internal/product/domain"
    appdb "r2-challenge/pkg/db"
    "r2-challenge/pkg/observability"
    "github.com/google/uuid"
)

type dbProductRepository struct {
    db     *gorm.DB
    tracer observability.Tracer
}

func NewDBRepository(database *appdb.Database, t observability.Tracer) (ProductRepository, error) {
    return &dbProductRepository{db: database.DB, tracer: t}, nil
}

func (r *dbProductRepository) Save(ctx context.Context, product domain.Product) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Save")
    defer span.End()

    now := time.Now().UTC()
    if product.ID == "" {
        product.ID = uuid.NewString()
    }
    product.CreatedAt = now
    product.UpdatedAt = now

    if err := r.db.WithContext(ctx).Table("products").Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "created_at"}, {Name: "updated_at"}}}).Create(&product).Error; err != nil {
        span.RecordError(err)
        return domain.Product{}, err
    }

    return product, nil
}

func (r *dbProductRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Update")
    defer span.End()

    product.UpdatedAt = time.Now().UTC()
    tx := r.db.WithContext(ctx).Table("products").Where("id = ?", product.ID).Updates(map[string]any{
        "name":        product.Name,
        "description": product.Description,
        "category":    product.Category,
        "price_cents": product.PriceCents,
        "inventory":   product.Inventory,
        "updated_at":  product.UpdatedAt,
    })
    if tx.Error != nil {
        span.RecordError(tx.Error)
        return domain.Product{}, tx.Error
    }
    if tx.RowsAffected == 0 {
        span.RecordError(gorm.ErrRecordNotFound)
        return domain.Product{}, gorm.ErrRecordNotFound
    }

    return r.GetByID(ctx, product.ID)
}

func (r *dbProductRepository) Delete(ctx context.Context, productID string) error {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.Delete")
    defer span.End()

    tx := r.db.WithContext(ctx).Table("products").Where("id = ?", productID).Delete(nil)
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

func (r *dbProductRepository) GetByID(ctx context.Context, productID string) (domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.GetByID")
    defer span.End()

    var product domain.Product
    err := r.db.WithContext(ctx).Table("products").Where("id = ?", productID).First(&product).Error
    if err != nil {
        span.RecordError(err)
        return domain.Product{}, err
    }

    return product, nil
}

func (r *dbProductRepository) List(ctx context.Context, filter ProductFilter) ([]domain.Product, error) {
    ctx, span := r.tracer.StartSpan(ctx, "ProductRepository.List")
    defer span.End()

    var list []domain.Product
    q := r.db.WithContext(ctx).Table("products")

    if filter.Category != "" {
        q = q.Where("category = ?", filter.Category)
    }
    if filter.Name != "" {
        q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
    }

    switch filter.SortBy {
    case "price":
        q = q.Order("price_cents " + order(filter.SortDesc))
    case "created_at":
        q = q.Order("created_at " + order(filter.SortDesc))
    default:
        q = q.Order("name " + order(filter.SortDesc))
    }

    if filter.Limit > 0 {
        q = q.Limit(filter.Limit)
    }
    if filter.Offset > 0 {
        q = q.Offset(filter.Offset)
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


