package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "microservice-template/common/api"
	commonerrors "microservice-template/common/errors"

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

func (s *repository) GetOrder(ctx context.Context, orderId *pb.OrderId) (*Order, error) {
	var order Order
	query := `
	SELECT 
		id,
		status,
		payment_link
	FROM orders
	WHERE id = $1
	`

	err := s.DB.Get(&order, query, orderId.ID)

	if err != nil {
		// no orders found
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No orders found.")
			return nil, commonerrors.ErrNoItemFound
		}

		fmt.Println("General Error when getting order of id %s, error: %s", orderId.ID, err)

		// general error
		return nil, err
	}

	fmt.Printf("Got order when querying for status %+v\n", order)

	return &order, nil
}

func (s repository) CreateOrder(ctx context.Context, order Order) (uuid.UUID, error) {
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

func (s *repository) UpdateOrderStatus(ctx context.Context, req *UpdateOrderStatusReq) error {

	fmt.Printf("\nUpdateOrderStatus repo: \nid: %s, status: %d\n\n", req.ID, req.Status)

	query := `
	UPDATE orders 
	SET 
		status = COALESCE(:status, status)
	WHERE id = :id
	`

	_, err := s.DB.NamedExecContext(ctx, query, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *repository) UpdateOrderPaymentLink(ctx context.Context, req *pb.OrderPaymentUpdateRequest) error {
	// temporary struct for db column mapping
	orderStruct := struct {
		ID          string `db:"id"`
		PaymentLink string `db:"payment_link"`
	}{
		ID:          req.ID,
		PaymentLink: req.PaymentLink,
	}

	query := `
	UPDATE orders 
	SET 
		payment_link = COALESCE(:payment_link, payment_link)
	WHERE id = :id
	`

	_, err := s.DB.NamedExecContext(ctx, query, orderStruct)
	if err != nil {
		return err
	}

	return nil
}

// transaction version
func (s *repository) CreateOrderTx(ctx context.Context, tx *sqlx.Tx, order Order) (uuid.UUID, error) {
	query := `
	INSERT INTO orders (customer_id, status)
	VALUES(:customer_id, :status)
	RETURNING id
	`

	// NOTE: sqlx missing NamedQueryContext method
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
