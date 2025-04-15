package app

import (
	"context"
	"fmt"
	"log"
	_ "solution/docs"
	v1 "solution/internal/api/handlers/v1"
	"solution/internal/api/middleware"
	"solution/internal/config"
	"solution/internal/database/mocks"
	"solution/internal/repository"
	"solution/internal/service"
	"solution/pkg/connections/postgres"
	"solution/pkg/metric"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type App struct {
	Config  *config.Config
	Fiber   *fiber.App
	DB      *postgres.DB
	MinIO   *minio.Client
	Metrics *metric.PromMetrics
}

func NewApp(ctx context.Context) *App {
	var app App
	app.Config = config.New()
	app.init(ctx)
	return &app
}

func (app *App) Start() {
	app.Fiber.Listen(app.Config.ServerAddr)
}

func (app *App) init(ctx context.Context) {
	app.connectDB(ctx)
	err := app.initMinIO(ctx)
	if err != nil {
		log.Println("failed to init minio:", err)
	}
	app.connectMetrics()
	app.engineSetup(ctx)
	app.handlersSetup()
}

func (app *App) connectDB(ctx context.Context) {
	app.DB = postgres.New(ctx, app.Config)
	app.DB.Migrate(ctx)
	app.DB.Pool()
	err := app.DB.Populate(ctx)
	if err != nil {
		fmt.Printf("failed to populate database: %v\n", err)
	}
}

func (app *App) initMinIO(ctx context.Context) error {
	cfg := app.Config
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKeyID, cfg.MinIOSecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	exists, err := minioClient.BucketExists(ctx, cfg.MinIOBucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinIOBucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return err
		}

		// Allow anonymous read-only access
		policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.MinIOBucketName)

		err = minioClient.SetBucketPolicy(ctx, cfg.MinIOBucketName, policy)
		if err != nil {
			return fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	app.MinIO = minioClient
	return nil
}

func (app *App) connectMetrics() {
	app.Metrics = &metric.PromMetrics{}
	app.Metrics.TotalRegistered = promauto.NewCounter(prometheus.CounterOpts{Name: "users_total_registered"})
	app.Metrics.TotalSeen = promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_seen"})
	app.Metrics.TotalLikedFilms = promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_liked"})
	app.Metrics.TotalDislikedFilms = promauto.NewCounter(prometheus.CounterOpts{Name: "films_total_disliked"})
	app.Metrics.TotalCouches = promauto.NewCounter(prometheus.CounterOpts{Name: "couches_total_created"})
}

func (app *App) engineSetup(ctx context.Context) {
	app.Fiber = fiber.New(fiber.Config{
		Prefork: false,
	})
	app.Fiber.Use(recover.New())
	app.Fiber.Use(logger.New())
	app.Fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
	}))
	app.Fiber.Use(middleware.CustomContext(ctx))
}

func (app *App) CloseConnections(ctx context.Context) {
	app.DB.Close(ctx)
}

func (app *App) handlersSetup() {
	// setting up prometheus
	app.Fiber.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	// route groups
	apiV1 := app.Fiber.Group("/api/v1")

	// api global routes
	apiV1.Get("docs/*", swagger.HandlerDefault)

	// global repos

	// global services
	authService := service.NewAuthService(app.Config)
	validatorService := service.NewValidatorService()

	// repos
	userRepo := repository.NewUserRepository(app.DB)
	cinemaRepo := repository.NewCinemaRepository(app.DB)
	couchRepo := repository.NewCouchRepository(app.DB)
	fileRepo := repository.NewFileRepo(app.MinIO, app.Config.MinIOBucketName)

	// services
	userService := service.NewUserService(app.Metrics, userRepo, authService)
	cinemaService := service.NewCinemaService(cinemaRepo)
	couchService := service.NewCouchService(app.Metrics, couchRepo)
	fileService := service.NewFileService(fileRepo, cinemaRepo, app.Config.MinIOPublicHost)

	// load pics to films
	picMocks := mocks.NewPicMocks(app.DB, fileService)
	picMocks.SetRandomPicsForCinemas()

	// handlers
	userHandler := v1.NewUserHandler(userService, authService, validatorService)

	cinemaHandler := v1.NewCinemaHandler(cinemaService, validatorService, authService, app.Config.AdminName, fileService)

	couchHandler := v1.NewCouchHandler(validatorService, authService, couchService, userService)

	oauthHandler := v1.NewOuathHandler(userService, authService)

	adminHandler := v1.NewAdminHandler(cinemaService, validatorService, authService)
	// handlers setup
	userHandler.Setup(apiV1, app.Config.SecretKey)
	cinemaHandler.Setup(apiV1, app.Config.SecretKey)

	cinemaHandler.Setup(apiV1, app.Config.SecretKey)
	couchHandler.Setup(apiV1, app.Config.SecretKey)
	adminHandler.Setup(apiV1, app.Config.SecretKey)
	oauthHandler.Setup(apiV1)
}
