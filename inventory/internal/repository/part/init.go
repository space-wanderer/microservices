package part

import (
	"sync"
	"time"

	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		parts: createSampleParts(),
	}
}

func createSampleParts() map[string]*repoModel.Part {
	parts := make(map[string]*repoModel.Part)

	// Двигатели
	parts["550e8400-e29b-41d4-a716-446655440001"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:          "Ионный двигатель X-2000",
		Description:   "Высокоэффективный ионный двигатель для межпланетных полетов",
		Price:         150000.0,
		StockQuantity: 5,
		Category:      repoModel.CategoryEngine,
		Dimensions: &repoModel.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 60.0,
			Weight: 250.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "КосмоТех",
			Country: "Россия",
			Website: "https://cosmotech.ru",
		},
		Tags:      []string{"ионный", "двигатель", "межпланетный", "высокоэффективный"},
		CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	parts["550e8400-e29b-41d4-a716-446655440002"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:          "Плазменный двигатель P-500",
		Description:   "Мощный плазменный двигатель для тяжелых грузов",
		Price:         200000.0,
		StockQuantity: 3,
		Category:      repoModel.CategoryEngine,
		Dimensions: &repoModel.Dimensions{
			Length: 150.0,
			Width:  100.0,
			Height: 80.0,
			Weight: 400.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "StarTech Industries",
			Country: "США",
			Website: "https://startech.com",
		},
		Tags:      []string{"плазменный", "двигатель", "тяжелый", "грузовой"},
		CreatedAt: time.Now().Add(-45 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Топливо
	parts["550e8400-e29b-41d4-a716-446655440003"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440003",
		Name:          "Криогенное топливо H2-O2",
		Description:   "Высокоэнергетическое криогенное топливо для ракетных двигателей",
		Price:         50000.0,
		StockQuantity: 20,
		Category:      repoModel.CategoryFuel,
		Dimensions: &repoModel.Dimensions{
			Length: 200.0,
			Width:  100.0,
			Height: 100.0,
			Weight: 1500.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "КриоТопливо",
			Country: "Россия",
			Website: "https://cryofuel.ru",
		},
		Tags:      []string{"криогенное", "топливо", "водород", "кислород"},
		CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	parts["550e8400-e29b-41d4-a716-446655440004"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440004",
		Name:          "Ядерное топливо U-235",
		Description:   "Обогащенный уран для ядерных реакторов",
		Price:         300000.0,
		StockQuantity: 2,
		Category:      repoModel.CategoryFuel,
		Dimensions: &repoModel.Dimensions{
			Length: 50.0,
			Width:  30.0,
			Height: 30.0,
			Weight: 100.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "АтомЭнерго",
			Country: "Россия",
			Website: "https://atomenergo.ru",
		},
		Tags:      []string{"ядерное", "топливо", "уран", "реактор"},
		CreatedAt: time.Now().Add(-90 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Иллюминаторы
	parts["550e8400-e29b-41d4-a716-446655440005"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440005",
		Name:          "Кварцевое окно QW-100",
		Description:   "Прозрачное кварцевое окно для космических кораблей",
		Price:         25000.0,
		StockQuantity: 15,
		Category:      repoModel.CategoryPorthole,
		Dimensions: &repoModel.Dimensions{
			Length: 100.0,
			Width:  100.0,
			Height: 10.0,
			Weight: 50.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "КварцТех",
			Country: "Россия",
			Website: "https://quartztech.ru",
		},
		Tags:      []string{"кварцевое", "окно", "прозрачное", "космическое"},
		CreatedAt: time.Now().Add(-20 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	parts["550e8400-e29b-41d4-a716-446655440006"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440006",
		Name:          "Бронированное окно BW-200",
		Description:   "Защищенное окно с многослойным покрытием",
		Price:         40000.0,
		StockQuantity: 8,
		Category:      repoModel.CategoryPorthole,
		Dimensions: &repoModel.Dimensions{
			Length: 120.0,
			Width:  120.0,
			Height: 15.0,
			Weight: 80.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "ArmorGlass",
			Country: "Германия",
			Website: "https://armorglass.de",
		},
		Tags:      []string{"бронированное", "окно", "защищенное", "многослойное"},
		CreatedAt: time.Now().Add(-15 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Крылья
	parts["550e8400-e29b-41d4-a716-446655440007"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440007",
		Name:          "Солнечная панель SP-500",
		Description:   "Высокоэффективная солнечная панель для космических станций",
		Price:         75000.0,
		StockQuantity: 12,
		Category:      repoModel.CategoryWing,
		Dimensions: &repoModel.Dimensions{
			Length: 500.0,
			Width:  200.0,
			Height: 5.0,
			Weight: 300.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "СолнТех",
			Country: "Россия",
			Website: "https://solntech.ru",
		},
		Tags:      []string{"солнечная", "панель", "энергия", "космическая"},
		CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	parts["550e8400-e29b-41d4-a716-446655440008"] = &repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440008",
		Name:          "Аэродинамическое крыло AW-300",
		Description:   "Легкое аэродинамическое крыло для атмосферных полетов",
		Price:         60000.0,
		StockQuantity: 10,
		Category:      repoModel.CategoryWing,
		Dimensions: &repoModel.Dimensions{
			Length: 300.0,
			Width:  150.0,
			Height: 20.0,
			Weight: 200.0,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "AeroDynamics",
			Country: "Франция",
			Website: "https://aerodynamics.fr",
		},
		Tags:      []string{"аэродинамическое", "крыло", "легкое", "атмосферное"},
		CreatedAt: time.Now().Add(-35 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	return parts
}
