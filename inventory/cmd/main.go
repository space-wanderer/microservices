package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpsPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "NotFound", req.GetUuid())
	}
	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Если фильтр не указан или все поля пусты - возвращаем все детали
	if req.GetFilter() == nil || isEmptyFilter(req.GetFilter()) {
		parts := make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	filter := req.GetFilter()
	var filteredParts []*inventoryV1.Part

	// Проходим по всем деталям и применяем фильтры
	for _, part := range s.parts {
		if matchesFilter(part, filter) {
			filteredParts = append(filteredParts, part)
		}
	}

	return &inventoryV1.ListPartsResponse{Parts: filteredParts}, nil
}

// isEmptyFilter проверяет, пуст ли фильтр
func isEmptyFilter(filter *inventoryV1.PartsFilter) bool {
	return len(filter.GetUuids()) == 0 &&
		len(filter.GetNames()) == 0 &&
		len(filter.GetCategories()) == 0 &&
		len(filter.GetManufacturerCountries()) == 0 &&
		len(filter.GetTags()) == 0
}

// matchesFilter проверяет, соответствует ли деталь всем условиям фильтра
func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	// Фильтр по UUID (логическое ИЛИ)
	if len(filter.GetUuids()) > 0 {
		uuidMatch := false
		for _, uuid := range filter.GetUuids() {
			if part.GetUuid() == uuid {
				uuidMatch = true
				break
			}
		}
		if !uuidMatch {
			return false
		}
	}

	// Фильтр по именам (логическое ИЛИ)
	if len(filter.GetNames()) > 0 {
		nameMatch := false
		for _, name := range filter.GetNames() {
			if part.GetName() == name {
				nameMatch = true
				break
			}
		}
		if !nameMatch {
			return false
		}
	}

	// Фильтр по категориям (логическое ИЛИ)
	if len(filter.GetCategories()) > 0 {
		categoryMatch := false
		for _, category := range filter.GetCategories() {
			if part.GetCategory() == category {
				categoryMatch = true
				break
			}
		}
		if !categoryMatch {
			return false
		}
	}

	// Фильтр по странам производителей (логическое ИЛИ)
	if len(filter.GetManufacturerCountries()) > 0 {
		countryMatch := false
		if part.GetManufacturer() != nil {
			for _, country := range filter.GetManufacturerCountries() {
				if part.GetManufacturer().GetCountry() == country {
					countryMatch = true
					break
				}
			}
		}
		if !countryMatch {
			return false
		}
	}

	// Фильтр по тегам (логическое ИЛИ)
	if len(filter.GetTags()) > 0 {
		tagMatch := false
		for _, filterTag := range filter.GetTags() {
			for _, partTag := range part.GetTags() {
				if partTag == filterTag {
					tagMatch = true
					break
				}
			}
			if tagMatch {
				break
			}
		}
		if !tagMatch {
			return false
		}
	}

	return true
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpsPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS inventory listening on %s", grpsPort)
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("gRPC server stopped")
}
