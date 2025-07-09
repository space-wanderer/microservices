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

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int32
	Category      inventoryV1.Category
	Dimensions    *inventoryV1.Dimensions
	Manufacturer  *inventoryV1.Manufacturer
	Tags          []string
	CreatedAt     *timestamppb.Timestamp
	UpdatedAt     *timestamppb.Timestamp
}

type Filter struct {
	Uuids                 []string
	Names                 []string
	Categories            []inventoryV1.Category
	ManufacturerCountries []string
	Tags                  []string
}

type InventoryStorage interface {
	Part(uuid string) (*Part, error)
	Parts(filter *Filter) ([]*Part, error)
}

type InventoryStorageInMem struct {
	mu    sync.RWMutex
	parts map[string]*Part
}

// Part - реализация метода интерфейса для получения детали по UUID
func (s *InventoryStorageInMem) Part(uuid string) (*Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[uuid]
	if !ok {
		return nil, fmt.Errorf("part not found: %s", uuid)
	}
	return part, nil
}

// Parts - реализация метода интерфейса для получения списка деталей с фильтрацией
func (s *InventoryStorageInMem) Parts(filter *Filter) ([]*Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if filter == nil || isEmptyFilter(filter) {
		parts := make([]*Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return parts, nil
	}

	var filteredParts []*Part

	for _, part := range s.parts {
		if matchesFilter(part, filter) {
			filteredParts = append(filteredParts, part)
		}
	}

	return filteredParts, nil
}

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage InventoryStorage
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := s.storage.Part(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "NotFound: %s", req.GetUuid())
	}

	grpcPart := convertPartToGRPC(part)
	return &inventoryV1.GetPartResponse{
		Part: grpcPart,
	}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := convertFilterFromGRPC(req.GetFilter())

	parts, err := s.storage.Parts(filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get parts: %v", err)
	}

	grpcParts := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		grpcParts[i] = convertPartToGRPC(part)
	}

	return &inventoryV1.ListPartsResponse{Parts: grpcParts}, nil
}

// convertPartToGRPC - конвертация внутренней модели Part в gRPC модель
func convertPartToGRPC(part *Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int64(part.StockQuantity),
		Category:      part.Category,
		Dimensions:    part.Dimensions,
		Manufacturer:  part.Manufacturer,
		Tags:          part.Tags,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// convertFilterFromGRPC - конвертация gRPC фильтра во внутренний фильтр
func convertFilterFromGRPC(grpcFilter *inventoryV1.PartsFilter) *Filter {
	if grpcFilter == nil {
		return &Filter{}
	}
	return &Filter{
		Uuids:                 grpcFilter.GetUuids(),
		Names:                 grpcFilter.GetNames(),
		Categories:            grpcFilter.GetCategories(),
		ManufacturerCountries: grpcFilter.GetManufacturerCountries(),
		Tags:                  grpcFilter.GetTags(),
	}
}

// isEmptyFilter проверяет, пуст ли фильтр
func isEmptyFilter(filter *Filter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

// matchesUuidFilter проверяет соответствие UUID
func matchesUuidFilter(part *Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	for _, uuid := range uuids {
		if part.Uuid == uuid {
			return true
		}
	}
	return false
}

// matchesNameFilter проверяет соответствие имени
func matchesNameFilter(part *Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if part.Name == name {
			return true
		}
	}
	return false
}

// matchesCategoryFilter проверяет соответствие категории
func matchesCategoryFilter(part *Part, categories []inventoryV1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, category := range categories {
		if part.Category == category {
			return true
		}
	}
	return false
}

// matchesCountryFilter проверяет соответствие страны производителя
func matchesCountryFilter(part *Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	for _, country := range countries {
		if part.Manufacturer.GetCountry() == country {
			return true
		}
	}
	return false
}

// matchesTagFilter проверяет соответствие тегов
func matchesTagFilter(part *Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, filterTag := range tags {
		for _, partTag := range part.Tags {
			if partTag == filterTag {
				return true
			}
		}
	}
	return false
}

// matchesFilter проверяет, соответствует ли деталь всем условиям фильтра
func matchesFilter(part *Part, filter *Filter) bool {
	return matchesUuidFilter(part, filter.Uuids) &&
		matchesNameFilter(part, filter.Names) &&
		matchesCategoryFilter(part, filter.Categories) &&
		matchesCountryFilter(part, filter.ManufacturerCountries) &&
		matchesTagFilter(part, filter.Tags)
}

// createSampleParts создает список деталей для тестирования
func createSampleParts() map[string]*Part {
	parts := make(map[string]*Part)

	// Двигатели
	parts["550e8400-e29b-41d4-a716-446655440001"] = &Part{
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

	parts["550e8400-e29b-41d4-a716-446655440002"] = &Part{
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
	parts["550e8400-e29b-41d4-a716-446655440003"] = &Part{
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

	parts["550e8400-e29b-41d4-a716-446655440004"] = &Part{
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
	parts["550e8400-e29b-41d4-a716-446655440005"] = &Part{
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

	parts["550e8400-e29b-41d4-a716-446655440006"] = &Part{
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
	parts["550e8400-e29b-41d4-a716-446655440007"] = &Part{
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

	parts["550e8400-e29b-41d4-a716-446655440008"] = &Part{
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

	// Создаем хранилище в памяти
	storage := &InventoryStorageInMem{
		parts: createSampleParts(),
	}

	// Создаем сервис с хранилищем через агрегацию
	service := &inventoryService{
		storage: storage,
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
