package main

import (
	"os"
	"os/signal"
	"product-service/internal/adapters"
	"product-service/internal/adapters/database"
	"product-service/internal/adapters/database/migrations"
	"product-service/internal/adapters/database/seeds"
	"product-service/internal/app/routers"
	"product-service/internal/domain/models"
	"product-service/internal/domain/repositories"
	"product-service/internal/domain/services"
	server "product-service/internal/grpc"
	"product-service/pkg/config"
	"product-service/pkg/constants"
	"product-service/pkg/logger"
	"product-service/pkg/utils"
	"product-service/pkg/workers"
	"sync"
	"syscall"
	"time"

	// product
	pbProduct "product-service/proto/product"
)

// @BasePath /api/v1
func main() {
	// init timezone and db
	initDB()
	defer cleanUp()
	// Set up shutdownCh and wg
	shutdownCh := make(chan struct{})
	var wg sync.WaitGroup

	/* Services */
	minIOService := adapters.NewMinIOService()
	kafkaService := adapters.NewKafkaService()
	// defer kafkaService.Close()
	openSearchService := adapters.NewOpenSearchService()

	/* init Kafka and Opensearch */
	initKafkaTopics(kafkaService)
	initOpenSearchIndexes(openSearchService)

	/* service */
	productVariantService := services.NewProductVariantService(minIOService, kafkaService, openSearchService)
	productService := services.NewProductService(minIOService, kafkaService, openSearchService, productVariantService)
	categoryService := services.NewCategoryService()
	manufacturerService := services.NewManufacturerService()
	tagService := services.NewTagService()
	brandService := services.NewBrandService()
	resourceService := services.NewResourceService(kafkaService, minIOService)

	// setup router
	router := routers.SetupRoute(productService, productVariantService, categoryService, manufacturerService, tagService, brandService, resourceService)

	// start kafka consumers
	go workers.StartNewProductKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartDeleteProductKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartUpdateProductKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartNewProductVariantKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartDeleteProductVariantKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartUpdateProductVariantKafkaConsumer(kafkaService, openSearchService, shutdownCh, &wg)
	go workers.StartProcessImageConsumer(kafkaService, minIOService, resourceService, shutdownCh, &wg)

	/* ingest initial seed data to opensearch documents, need to run at once */
	ingestSeedDataToOpensearch(kafkaService, openSearchService)

	// Start the Gin server concurrently in a Goroutine
	serverErrCh := make(chan error)
	go func() {
		serverErrCh <- router.Run(config.ServerConfig())
	}()
	// start grpc server
	server.StartServer(productService, productVariantService, categoryService, manufacturerService, tagService, brandService, resourceService)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Wait for a signal or a server error
	select {
	case sig := <-sigCh:
		logger.Infof("Received signal: %v. Shutting down...", sig)
	case err := <-serverErrCh:
		logger.Infof("Server error: %v. Shutting down...", err)
	}

	close(shutdownCh) // signal all goroutines to stop
	wg.Wait()         // wait for all goroutines to stop
}

func initDB() {
	//set timezone
	os.Setenv("SERVER_TIMEZONE", "Asia/Tokyo")
	loc, _ := time.LoadLocation(os.Getenv("SERVER_TIMEZONE"))
	time.Local = loc
	// setup db
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	// run db migration
	migrations.Migrate()
	// Check if seed data exists
	if !seeds.IsSeedDataExists() {
		// Run db seed
		seeds.SeedData()
	}
}

func cleanUp() {

}

func ingestSeedDataToOpensearch(kafkaService *adapters.KafkaService, openSearchService *adapters.OpenSearchService) {
	if !openSearchService.IndexExists(constants.OPENSEARCH_PRODUCT_INDEX) {
		logger.Errorf("Error no index for product")
		return
	}

	documentCount, _ := openSearchService.DocumentCount(constants.OPENSEARCH_PRODUCT_INDEX)
	logger.Infof("Opensearch product document count: %d", documentCount)
	if documentCount > 0 {
		return
	}

	productRepo := &repositories.ProductRepository{}
	products, err := productRepo.GetAll()
	if err != nil {
		logger.Errorf("Error getting products from db: %s", err)
		return
	}
	productVariantRepo := &repositories.ProductVariantRepository{}
	productVariants, err := productVariantRepo.GetAll()
	if err != nil {
		logger.Errorf("Error getting product variants from db: %s", err)
		return
	}

	// Create a channel to receive errors from goroutines
	errorChannel := make(chan error, len(products)+len(productVariants))

	// Use a wait group to wait for all goroutines to finish
	var ingestSeedWg sync.WaitGroup

	// remove all indexed documents
	openSearchService.DeleteAllDocuments(constants.OPENSEARCH_PRODUCT_INDEX)
	openSearchService.DeleteAllDocuments(constants.OPENSEARCH_PRODUCT_VARIANT_INDEX)

	for _, product := range products {
		ingestSeedWg.Add(1)

		go func(p models.Product) {
			defer ingestSeedWg.Done()
			// convert models.Product to pbProduct.Product
			var product pbProduct.Product
			err := utils.CopyProductModelToProto(&product, &p)
			if err != nil {
				errorChannel <- err
				return
			}
			productString, err := utils.ConvertProtoToByte(&product)
			if err != nil {
				errorChannel <- err
				return
			}
			err = kafkaService.ProduceMessageWithRetry(constants.KAFKA_NEW_PRODUCT_TOPIC, productString, constants.KAFKA_PRODUCER_MAX_RETRIES)
			// err = openSearchService.IndexProduct(productString)
			if err != nil {
				errorChannel <- err
			}
		}(product)
	}

	for _, productVariant := range productVariants {
		ingestSeedWg.Add(1)

		go func(p models.ProductVariant) {
			defer ingestSeedWg.Done()
			// convert models.ProductVariant to pbProduct.ProductVariantWithProduct
			var productVariant pbProduct.ProductVariantWithProduct
			err := utils.CopyProductVariantModelToProto(&productVariant, &p)
			if err != nil {
				errorChannel <- err
				return
			}
			productVariantString, err := utils.ConvertProtoToByte(&productVariant)
			if err != nil {
				errorChannel <- err
				return
			}
			err = kafkaService.ProduceMessageWithRetry(constants.KAFKA_NEW_PRODUCT_VARIANT_TOPIC, productVariantString, constants.KAFKA_PRODUCER_MAX_RETRIES)
			// err = openSearchService.IndexProductVariant(productVariantString)
			if err != nil {
				errorChannel <- err
			}
		}(productVariant)
	}

	// Wait for all goroutines to finish
	ingestSeedWg.Wait()

	// Close the error channel after all goroutines are done
	close(errorChannel)

	// Check for errors from goroutines
	for err := range errorChannel {
		logger.Errorf("Error producing Kafka message: %s", err)
	}
}

func initKafkaTopics(kafkaService *adapters.KafkaService) {
	if !kafkaService.TopicExists(constants.KAFKA_NEW_PRODUCT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_NEW_PRODUCT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
	if !kafkaService.TopicExists(constants.KAFKA_DELETE_PRODUCT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_DELETE_PRODUCT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
	if !kafkaService.TopicExists(constants.KAFKA_UPDATE_PRODUCT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_UPDATE_PRODUCT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
	if !kafkaService.TopicExists(constants.KAFKA_NEW_PRODUCT_VARIANT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_NEW_PRODUCT_VARIANT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
	if !kafkaService.TopicExists(constants.KAFKA_DELETE_PRODUCT_VARIANT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_DELETE_PRODUCT_VARIANT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
	if !kafkaService.TopicExists(constants.KAFKA_UPDATE_PRODUCT_VARIANT_TOPIC) {
		kafkaService.CreateTopic(constants.KAFKA_UPDATE_PRODUCT_VARIANT_TOPIC, constants.KAFKA_PRODUCER_PARTITIONS, constants.KAFKA_PRODUCER_REPLICATION_FACTOR)
	}
}

func initOpenSearchIndexes(openSearchService *adapters.OpenSearchService) {
	if !openSearchService.IndexExists(constants.OPENSEARCH_PRODUCT_INDEX) {
		err := openSearchService.CreateProductsIndex()
		if err != nil {
			logger.Fatalf("Error creating products index on Opensearch: %s", err)
		}
	}
	if !openSearchService.IndexExists(constants.OPENSEARCH_PRODUCT_VARIANT_INDEX) {
		err := openSearchService.CreateProductVariantsIndex()
		if err != nil {
			logger.Fatalf("Error creating product variants index on Opensearch: %s", err)
		}
	}
}
