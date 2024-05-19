package memory

import (
	"sync"

	"github.com/codepnw/go-ddd/aggregate"
	"github.com/codepnw/go-ddd/domain/product"
	"github.com/google/uuid"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]aggregate.Product
	sync.Mutex
}

func New() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]aggregate.Product),
	}
}

func (pr *MemoryProductRepository) GetAll() ([]aggregate.Product, error) {
	var products []aggregate.Product

	for _, product := range pr.products {
		products = append(products, product)
	}
	return products, nil
}

func (pr *MemoryProductRepository) GetByID(id uuid.UUID) (aggregate.Product, error) {
	if product, ok := pr.products[id]; ok {
		return product, nil
	}
	return aggregate.Product{}, product.ErrProductNotFound
}

func (pr *MemoryProductRepository) Add(newprod aggregate.Product) error {
	pr.Lock()
	defer pr.Unlock()

	if _, ok := pr.products[newprod.GetID()]; ok {
		return product.ErrProductAlreadyExists
	}

	pr.products[newprod.GetID()] = newprod

	return nil
}

func (pr *MemoryProductRepository) Update(update aggregate.Product) error {
	pr.Lock()
	defer pr.Unlock()

	if _, ok := pr.products[update.GetID()]; !ok {
		return product.ErrProductNotFound
	}

	pr.products[update.GetID()] = update
	return nil
}

func (pr *MemoryProductRepository) Delete(id uuid.UUID) error {
	pr.Lock()
	defer pr.Unlock()

	if _, ok := pr.products[id]; !ok {
		return product.ErrProductNotFound
	}
	delete(pr.products, id)

	return nil
}
