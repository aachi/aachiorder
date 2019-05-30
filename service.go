package aachiorder

import (
	"context"
	"errors"
)

type Order struct {
	ID           string      `json:"id,omitempty"`
	CustomerID   string      `json:"customer_id"`
	Status       string      `json:"status"`
	CreatedOn    int64       `json:"created_on,omitempty"`
	RestaurantId string      `json:"restaurant_id"`
	OrderItems   []OrderItem `json:"order_items,omitempty"`
}

// OrderItem represents items in an order
type OrderItem struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}
// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	Order order.Order
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Order order.Order `json:"order"`
	Err   error       `json:"error,omitempty"`
}

// ChangeStatusRequest holds the request parameters for the ChangeStatus method.
type ChangeStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ChangeStatusResponse holds the response values for the ChangeStatus method.
type ChangeStatusResponse struct {
	Err error `json:"error,omitempty"`
}

type Service interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, id string) (Order, error)
	ChangeStatus(ctx context.Context, id string, status string) error
}

type service struct {
	repository ordersvc.Repository
	logger     log.Logger
}

// NewService creates and returns a new Order service instance
func NewService(rep ordersvc.Repository, logger log.Logger) ordersvc.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Create makes an order
func (s *service) Create(ctx context.Context, order ordersvc.Order) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	order.ID = id
	order.Status = "Pending"
	order.CreatedOn = time.Now().Unix()

	if err := s.repository.CreateOrder(ctx, order); err != nil {
		level.Error(logger).Log("err", err)
		return "", ordersvc.ErrCmdRepository
	}
	return id, nil
}

// GetByID returns an order given by id
func (s *service) GetByID(ctx context.Context, id string) (ordersvc.Order, error) {
	logger := log.With(s.logger, "method", "GetByID")
	order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return order, ordersvc.ErrOrderNotFound
		}
		return order, ordersvc.ErrQueryRepository
	}
	return order, nil
}

// ChangeStatus changes the status of an order
func (s *service) ChangeStatus(ctx context.Context, id string, status string) error {
	logger := log.With(s.logger, "method", "ChangeStatus")
	if err := s.repository.ChangeOrderStatus(ctx, id, status); err != nil {
		level.Error(logger).Log("err", err)
		return ordersvc.ErrCmdRepository
	}
	return nil
}
