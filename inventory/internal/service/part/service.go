package part

import "github.com/space-wanderer/microservices/inventory/internal/repository"

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{inventoryRepository: inventoryRepository}
}
