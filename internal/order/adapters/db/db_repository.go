package db

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"r2-challenge/internal/order/domain"
	appdb "r2-challenge/pkg/db"
	"r2-challenge/pkg/observability"
	"github.com/google/uuid"
)

type dbOrderRepository struct {
	db     *gorm.DB
	tracer observability.Tracer
}

func NewDBRepository(database *appdb.Database, t observability.Tracer) (OrderRepository, error) {
	return &dbOrderRepository{db: database.DB, tracer: t}, nil
}

func (r *dbOrderRepository) Save(ctx context.Context, order domain.Order) (domain.Order, error) {
	ctx, span := r.tracer.StartSpan(ctx, "OrderRepository.Save")
	defer span.End()

	now := time.Now().UTC()
	if order.ID == "" { 
		order.ID = uuid.NewString()
	}
	order.CreatedAt = now
	order.UpdatedAt = now

	tx := r.db.WithContext(ctx).Begin()

	// prevent auto-saving associations (items)
	if err := tx.Omit(clause.Associations).Table("orders").Create(&order).Error; err != nil {
		span.RecordError(err)
		tx.Rollback()
		return domain.Order{}, err
	}

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		if order.Items[i].ID == "" { 
			order.Items[i].ID = uuid.NewString()
		}
	}

	if len(order.Items) > 0 {
		// Atomic inventory check and decrement per item
		for _, it := range order.Items {
			res := tx.Exec("UPDATE products SET inventory = inventory - ? WHERE id = ? AND inventory >= ?", it.Quantity, it.ProductID, it.Quantity)
			if res.Error != nil {
				span.RecordError(res.Error)
				tx.Rollback()
				return domain.Order{}, res.Error
			}
			if res.RowsAffected == 0 {
				tx.Rollback()
				return domain.Order{}, gorm.ErrInvalidData
			}
		}
		if err := tx.Table("order_items").Omit("id").Create(&order.Items).Error; err != nil {
			span.RecordError(err)
			tx.Rollback()
			return domain.Order{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}

	return order, nil
}

func (r *dbOrderRepository) UpdateStatus(ctx context.Context, id string, status string) (domain.Order, error) {
	ctx, span := r.tracer.StartSpan(ctx, "OrderRepository.UpdateStatus")
	defer span.End()

	if err := r.db.WithContext(ctx).Table("orders").Where("id = ?", id).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}).Error; err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}

	return r.GetByID(ctx, id)
}

func (r *dbOrderRepository) GetByID(ctx context.Context, orderID string) (domain.Order, error) {
	ctx, span := r.tracer.StartSpan(ctx, "OrderRepository.GetByID")
	defer span.End()

	var order domain.Order
	if err := r.db.WithContext(ctx).Table("orders").Where("id = ?", orderID).First(&order).Error; err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}
	var items []domain.OrderItem

	if err := r.db.WithContext(ctx).Table("order_items").Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}
	order.Items = items

	return order, nil
}

func (r *dbOrderRepository) ListByUser(ctx context.Context, userID string, filter OrderFilter) ([]domain.Order, error) {
	ctx, span := r.tracer.StartSpan(ctx, "OrderRepository.ListByUser")
	defer span.End()

	var orders []domain.Order
	query := r.db.WithContext(ctx).Table("orders").Where("user_id = ?", userID)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&orders).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	if len(orders) == 0 {
		return orders, nil
	}

	ids := make([]string, 0, len(orders))
	for _, ord := range orders {
		ids = append(ids, ord.ID)
	}

	var items []domain.OrderItem
	if err := r.db.WithContext(ctx).Table("order_items").Where("order_id IN ?", ids).Find(&items).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	itemsByOrder := make(map[string][]domain.OrderItem, len(orders))
	for _, it := range items {
		itemsByOrder[it.OrderID] = append(itemsByOrder[it.OrderID], it)
	}

	for i := range orders {
		orders[i].Items = itemsByOrder[orders[i].ID]
	}

	return orders, nil
}


