package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().GetName()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetDescription()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetPrice()).To(BeNumerically(">", 0))
			Expect(resp.GetPart().GetStockQuantity()).To(BeNumerically(">=", 0))
			Expect(resp.GetPart().GetCategory()).ToNot(BeNil())
			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
			Expect(resp.GetPart().GetManufacturer()).ToNot(BeNil())
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
		})

		It("должен возвращать ошибку для несуществующего UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: "non-existent-uuid",
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Вставляем несколько тестовых деталей
			_, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())
			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())
			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("должен успешно возвращать список всех деталей", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(resp.GetParts()).To(HaveLen(3))

			// Проверяем, что каждая деталь имеет необходимые поля
			for _, part := range resp.GetParts() {
				Expect(part.GetUuid()).ToNot(BeEmpty())
				Expect(part.GetName()).ToNot(BeEmpty())
				Expect(part.GetDescription()).ToNot(BeEmpty())
				Expect(part.GetPrice()).To(BeNumerically(">", 0))
				Expect(part.GetStockQuantity()).To(BeNumerically(">=", 0))
				Expect(part.GetCategory()).ToNot(BeNil())
				Expect(part.GetDimensions()).ToNot(BeNil())
				Expect(part.GetManufacturer()).ToNot(BeNil())
			}
		})

		It("должен фильтровать детали по категории", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())

			// Проверяем, что все возвращенные детали имеют категорию ENGINE
			for _, part := range resp.GetParts() {
				Expect(part.GetCategory()).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
			}
		})

		It("должен фильтровать детали по стране производителя", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					ManufacturerCountries: []string{"Россия"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())

			// Проверяем, что все возвращенные детали произведены в России
			for _, part := range resp.GetParts() {
				Expect(part.GetManufacturer().GetCountry()).To(Equal("Россия"))
			}
		})

		It("должен фильтровать детали по тегам", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Tags: []string{"двигатель"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())

			// Проверяем, что все возвращенные детали содержат тег "двигатель"
			for _, part := range resp.GetParts() {
				Expect(part.GetTags()).To(ContainElement("двигатель"))
			}
		})

		It("должен возвращать пустой список для несуществующих фильтров", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_UNSPECIFIED},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeEmpty())
		})
	})

	Describe("Полный жизненный цикл", func() {
		It("должен поддерживать получение и фильтрацию деталей", func() {
			// 1. Вставляем тестовую деталь с известными данными
			testPart := env.GetTestPartInfo()
			partUUID, err := env.InsertTestPartWithData(ctx, testPart)
			Expect(err).ToNot(HaveOccurred())

			// 2. Получаем деталь по UUID
			getResp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getResp.GetPart().Uuid).To(Equal(partUUID))
			Expect(getResp.GetPart().GetName()).To(Equal(testPart.GetName()))
			Expect(getResp.GetPart().GetDescription()).To(Equal(testPart.GetDescription()))
			Expect(getResp.GetPart().GetPrice()).To(Equal(testPart.GetPrice()))
			Expect(getResp.GetPart().GetStockQuantity()).To(Equal(testPart.GetStockQuantity()))
			Expect(getResp.GetPart().GetCategory()).To(Equal(testPart.GetCategory()))

			// 3. Получаем список деталей с фильтром по категории
			listResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что наша деталь есть в списке
			found := false
			for _, part := range listResp.GetParts() {
				if part.GetUuid() == partUUID {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// 4. Получаем список деталей с фильтром по стране производителя
			listByCountryResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					ManufacturerCountries: []string{"Россия"},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(listByCountryResp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что наша деталь есть в списке
			found = false
			for _, part := range listByCountryResp.GetParts() {
				if part.GetUuid() == partUUID {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})
	})
})
