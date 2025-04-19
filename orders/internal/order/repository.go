package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) OrderRepository {
	return &repository{
		DB: db,
	}
}

func (s *repository) CreateOrder(ctx context.Context, order Order) (uuid.UUID, error) {
	query := `
	INSERT INTO orders (customer_id, status)
	VALUES(:customer_id, :status)
	RETURNING id
	`

	rows, err := s.DB.NamedQueryContext(ctx, query, order)

	if err != nil {
		return uuid.Nil, err
	}

	defer rows.Close()

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

// transaction version
func (s *repository) CreateOrderTx(ctx context.Context, tx *sqlx.Tx, order Order) (uuid.UUID, error) {
	query := `
	INSERT INTO orders (customer_id, status)
	VALUES(:customer_id, :status)
	RETURNING id
	`

	rows, err := tx.NamedQuery(query, order)

	if err != nil {
		return uuid.Nil, err
	}

	defer rows.Close()

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
	INSERT INTO order_items(order_id, name, quantity, price_id)
	VALUES(:order_id, :name, :quantity, :price_id)
	`

	_, err := s.DB.NamedExec(query, item)

	if err != nil {
		return err
	}

	fmt.Println("Successfully created order item.")
	return nil
}

// transaction version
func (s *repository) CreateOrderItemTx(ctx context.Context, tx *sqlx.Tx, item OrderItem) error {
	query := `
	INSERT INTO order_items(order_id, name, quantity, price_id)
	VALUES(:order_id, :name, :quantity, :price_id)
	`

	_, err := tx.NamedExec(query, item)

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

	err := s.DB.Select(&orders, query)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
