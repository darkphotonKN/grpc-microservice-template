package payment

import (
	"context"
	pb "microservice-template/common/api"
	"microservice-template/payments/internal/order"
	"microservice-template/payments/processor/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	processor *inmem.InMemProcessor
	service   *service
}

// setup fields with respective fields - method name fixed
func (s *PaymentServiceTestSuite) SetupTest() {
	newInMemProcessor := inmem.NewInMemProcessor()
	s.processor = newInMemProcessor.(*inmem.InMemProcessor)

	newPaymentService := NewService(newInMemProcessor, "")
	s.service = newPaymentService
}

// happy path test
func (s *PaymentServiceTestSuite) TestCreatePaymentSuccess() {

	order := &pb.Order{
		ID:         "testID123",
		CustomerID: "testCustomer123",
		Status:     "testStatus",
		Items: []*pb.Item{
			&pb.Item{
				ID:       "123",
				Name:     "item 1",
				Quantity: 3,
				PriceID:  "123",
			},
			&pb.Item{
				ID:       "123",
				Name:     "item 1",
				Quantity: 3,
				PriceID:  "123",
			},
		},
	}

	paymentLink, err := s.service.CreatePayment(context.Background(), order)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), paymentLink)
	assert.Equal(s.T(), "test link", paymentLink)
}

// run tests
func TestPaymentService(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
