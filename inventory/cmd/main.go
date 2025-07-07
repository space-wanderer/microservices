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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
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
		return nil, status.Errorf(codes.NotFound, "NotFound: %s", req.GetUuid())
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

// matchesUuidFilter проверяет соответствие UUID
func matchesUuidFilter(part *inventoryV1.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	for _, uuid := range uuids {
		if part.GetUuid() == uuid {
			return true
		}
	}
	return false
}

// matchesNameFilter проверяет соответствие имени
func matchesNameFilter(part *inventoryV1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if part.GetName() == name {
			return true
		}
	}
	return false
}

// matchesCategoryFilter проверяет соответствие категории
func matchesCategoryFilter(part *inventoryV1.Part, categories []inventoryV1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, category := range categories {
		if part.GetCategory() == category {
			return true
		}
	}
	return false
}

// matchesCountryFilter проверяет соответствие страны производителя
func matchesCountryFilter(part *inventoryV1.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.GetManufacturer() == nil {
		return false
	}
	for _, country := range countries {
		if part.GetManufacturer().GetCountry() == country {
			return true
		}
	}
	return false
}

// matchesTagFilter проверяет соответствие тегов
func matchesTagFilter(part *inventoryV1.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, filterTag := range tags {
		for _, partTag := range part.GetTags() {
			if partTag == filterTag {
				return true
			}
		}
	}
	return false
}

// matchesFilter проверяет, соответствует ли деталь всем условиям фильтра
func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	return matchesUuidFilter(part, filter.GetUuids()) &&
		matchesNameFilter(part, filter.GetNames()) &&
		matchesCategoryFilter(part, filter.GetCategories()) &&
		matchesCountryFilter(part, filter.GetManufacturerCountries()) &&
		matchesTagFilter(part, filter.GetTags())
}

// createSampleParts создает список деталей для тестирования
func createSampleParts() map[string]*inventoryV1.Part {
	parts := make(map[string]*inventoryV1.Part)

	// Двигатели
	parts["550e8400-e29b-41d4-a716-446655440001"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440001",
		Name:          "Ионный двигатель X-2000",
		Description:   "Высокоэффективный ионный двигатель для межпланетных полетов",
		Price:         150000.0,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 60.0,
			Weight: 250.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "КосмоТех",
			Country: "Россия",
			Website: "https://cosmotech.ru",
		},
		Tags:      []string{"ионный", "двигатель", "межпланетный", "высокоэффективный"},
		CreatedAt: timestamppb.New(time.Now().Add(-30 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	parts["550e8400-e29b-41d4-a716-446655440002"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440002",
		Name:          "Плазменный двигатель P-500",
		Description:   "Мощный плазменный двигатель для тяжелых грузов",
		Price:         200000.0,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 150.0,
			Width:  100.0,
			Height: 80.0,
			Weight: 400.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "StarTech Industries",
			Country: "США",
			Website: "https://startech.com",
		},
		Tags:      []string{"плазменный", "двигатель", "тяжелый", "грузовой"},
		CreatedAt: timestamppb.New(time.Now().Add(-45 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Топливо
	parts["550e8400-e29b-41d4-a716-446655440003"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440003",
		Name:          "Криогенное топливо H2-O2",
		Description:   "Высокоэнергетическое криогенное топливо для ракетных двигателей",
		Price:         50000.0,
		StockQuantity: 20,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 200.0,
			Width:  100.0,
			Height: 100.0,
			Weight: 1500.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "КриоТопливо",
			Country: "Россия",
			Website: "https://cryofuel.ru",
		},
		Tags:      []string{"криогенное", "топливо", "водород", "кислород"},
		CreatedAt: timestamppb.New(time.Now().Add(-60 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	parts["550e8400-e29b-41d4-a716-446655440004"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440004",
		Name:          "Ядерное топливо U-235",
		Description:   "Обогащенный уран для ядерных реакторов",
		Price:         300000.0,
		StockQuantity: 2,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 50.0,
			Width:  30.0,
			Height: 30.0,
			Weight: 100.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "АтомЭнерго",
			Country: "Россия",
			Website: "https://atomenergo.ru",
		},
		Tags:      []string{"ядерное", "топливо", "уран", "реактор"},
		CreatedAt: timestamppb.New(time.Now().Add(-90 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Иллюминаторы
	parts["550e8400-e29b-41d4-a716-446655440005"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440005",
		Name:          "Кварцевое окно QW-100",
		Description:   "Прозрачное кварцевое окно для космических кораблей",
		Price:         25000.0,
		StockQuantity: 15,
		Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 100.0,
			Width:  100.0,
			Height: 10.0,
			Weight: 50.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "КварцТех",
			Country: "Россия",
			Website: "https://quartztech.ru",
		},
		Tags:      []string{"кварцевое", "окно", "прозрачное", "космическое"},
		CreatedAt: timestamppb.New(time.Now().Add(-20 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	parts["550e8400-e29b-41d4-a716-446655440006"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440006",
		Name:          "Бронированное окно BW-200",
		Description:   "Защищенное окно с многослойным покрытием",
		Price:         40000.0,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.0,
			Width:  120.0,
			Height: 15.0,
			Weight: 80.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "ArmorGlass",
			Country: "Германия",
			Website: "https://armorglass.de",
		},
		Tags:      []string{"бронированное", "окно", "защищенное", "многослойное"},
		CreatedAt: timestamppb.New(time.Now().Add(-15 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Крылья
	parts["550e8400-e29b-41d4-a716-446655440007"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440007",
		Name:          "Солнечная панель SP-500",
		Description:   "Высокоэффективная солнечная панель для космических станций",
		Price:         75000.0,
		StockQuantity: 12,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 500.0,
			Width:  200.0,
			Height: 5.0,
			Weight: 300.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "СолнТех",
			Country: "Россия",
			Website: "https://solntech.ru",
		},
		Tags:      []string{"солнечная", "панель", "энергия", "космическая"},
		CreatedAt: timestamppb.New(time.Now().Add(-25 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	parts["550e8400-e29b-41d4-a716-446655440008"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440008",
		Name:          "Аэродинамическое крыло AW-300",
		Description:   "Легкое аэродинамическое крыло для атмосферных полетов",
		Price:         60000.0,
		StockQuantity: 10,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 300.0,
			Width:  150.0,
			Height: 20.0,
			Weight: 200.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroDynamics",
			Country: "Франция",
			Website: "https://aerodynamics.fr",
		},
		Tags:      []string{"аэродинамическое", "крыло", "легкое", "атмосферное"},
		CreatedAt: timestamppb.New(time.Now().Add(-35 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	return parts
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
		parts: createSampleParts(),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS inventory listening on %d", grpsPort)
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
