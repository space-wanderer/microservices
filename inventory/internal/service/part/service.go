package part

import "github.com/space-wanderer/microservices/inventory/internal/service"

type Service struct {
	inventoryRepository service.PartService
}

func NewService(inventoryRepository service.PartService) *Service {
	return &Service{inventoryRepository: inventoryRepository}
}
