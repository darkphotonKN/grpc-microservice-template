package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) OrderRepository {
	return &repository{
		db: db,
	}
}

func (s *repository) CreateOrder(ctx context.Context, order Order) (uuid.UUID, error) {
	query := `
	INSERT INTO orders (customer_id, status)
	VALUES(:customer_id, :status)
	RETURNING id
	`

	rows, err := s.db.NamedQueryContext(ctx, query, order)
	defer rows.Close()

	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return uuid.Nil, err
		}
	}

	fmt.Printf("Successfully created order with id: %s\n", id)

	return id, nil
}

func (s *repository) CreateOrderItem(ctx context.Context, item OrderItem) error {
	query := `
	INSERT INTO order_items 
	VALUES(:order_id, :name, :quantity, :price_id)
	`

	_, err := s.db.NamedExec(query, item)

	if err != nil {
		return err
	}

	fmt.Println("Successfully created order item.")
	return nil
}

func (s *repository) GetAll(ctx context.Context) ([]*pb.Order, error) {
	var orders []*pb.Order

	query := `
	SELECT 
		id
		user_id
		product_id
		quantity
		price
		status
		created_at
		updated_at
	FROM orders
	`

	err := s.db.Select(&orders, query)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
