package services

import (
	"log"

	"github.com/codepnw/go-ddd/aggregate"
	"github.com/codepnw/go-ddd/domain/customer"
	"github.com/codepnw/go-ddd/domain/customer/memory"
	"github.com/codepnw/go-ddd/domain/product"
	prodmem "github.com/codepnw/go-ddd/domain/product/memory"
	"github.com/google/uuid"
)

type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	customers customer.CustomerRepository
	products  product.ProductRepository
}

func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}

	for _, cfg := range cfgs {
		err := cfg(os)

		if err != nil {
			return nil, err
		}
	}

	return os, nil
}

// applies a customer repository to the OrderService
func WithCostomerRepository(cr customer.CustomerRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

func WithMemoryCustomerRepository() OrderConfiguration {
	cr := memory.New()
	return WithCostomerRepository(cr)
}

func WithMemoryProductRepository(products []aggregate.Product) OrderConfiguration {
	return func(os *OrderService) error {
		pr := prodmem.New()

		for _, p := range products {
			if err := pr.Add(p); err != nil {
				return err
			}
		}
		return nil
	}
}

func (o *OrderService) CreateOrder(customerID uuid.UUID, productsIDs []uuid.UUID) (float64, error) {
	// fetch the customer
	c, err := o.customers.Get(customerID)
	if err != nil {
		return 0, err 
	}

	var products []aggregate.Product
	var total float64

	for _, id := range productsIDs {
		p, err := o.products.GetByID(id)

		if err != nil {
			return 0, err
		}

		products = append(products, p)
		total += p.GetPrice()
	}

	log.Printf("Customer: %s has orderd %d products", c.GetID(), len(products))
	return total, nil
}
