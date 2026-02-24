package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"caltrack/config"
	_ "caltrack/docs"
	"caltrack/handler/analyze"
	"caltrack/handler/auth"
	"caltrack/handler/middleware"
	"caltrack/handler/nutrition"
	"caltrack/handler/record"
	"caltrack/handler/user"
	gormPersistence "caltrack/infrastructure/persistence/gorm"
	infraService "caltrack/infrastructure/service"
	"caltrack/pkg/logger"
	"caltrack/usecase"
)

// @title CalTrack API
// @version 1.0
// @description カロリー管理アプリケーションのAPI
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// ロガー初期化
	logger.Init()

	// マイグレーションを実行
	if err := config.RunMigrations(); err != nil {
		logger.Error("マイグレーション失敗", "error", err.Error())
		panic(err)
	}

	// DB接続
	database, err := config.NewDatabase()
	if err != nil {
		logger.Error("DB接続失敗", "error", err.Error())
		panic(err)
	}

	// Geminiクライアント初期化
	geminiConfig, err := config.NewGeminiConfig()
	if err != nil {
		logger.Error("Geminiクライアント初期化失敗", "error", err.Error())
		panic(err)
	}
	defer geminiConfig.Close()

	// DI - Repository
	userRepo := gormPersistence.NewGormUserRepository(database.DB)
	sessionRepo := gormPersistence.NewGormSessionRepository(database.DB)
	recordRepo := gormPersistence.NewGormRecordRepository(database.DB)
	recordPfcRepo := gormPersistence.NewGormRecordPfcRepository(database.DB)
	adviceCacheRepo := gormPersistence.NewGormAdviceCacheRepository(database.DB)
	txManager := gormPersistence.NewGormTransactionManager(database.DB)

	// DI - Service
	imageAnalyzer := infraService.NewGeminiImageAnalyzer(geminiConfig.Client)
	pfcAnalyzer := infraService.NewGeminiPfcAnalyzer(geminiConfig.Client)
	pfcEstimator := infraService.NewGeminiPfcEstimator(geminiConfig.Client)

	// DI - Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, txManager)
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	recordUsecase := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, geminiConfig)
	analyzeUsecase := usecase.NewAnalyzeUsecase(imageAnalyzer, geminiConfig)
	nutritionUsecase := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, pfcAnalyzer, geminiConfig)

	// DI - Handler
	userHandler := user.NewUserHandler(userUsecase)
	authHandler := auth.NewAuthHandler(authUsecase)
	recordHandler := record.NewRecordHandler(recordUsecase)
	analyzeHandler := analyze.NewAnalyzeHandler(analyzeUsecase)
	nutritionHandler := nutrition.NewNutritionHandler(nutritionUsecase)

	// Setup router
	r := gin.Default()

	// ロガーミドルウェア追加
	r.Use(logger.Middleware())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetCORSAllowOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "CalTrack API is running"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.Register)
	}

	// Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// 認証が必要なルート
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware(authUsecase))
	{
		authenticated.GET("/users/profile", userHandler.GetProfile)
		authenticated.PATCH("/users/profile", userHandler.UpdateProfile)
		authenticated.POST("/records", recordHandler.Create)
		authenticated.GET("/records/today", recordHandler.GetToday)
		authenticated.GET("/statistics", recordHandler.GetStatistics)
		authenticated.POST("/analyze-image", analyzeHandler.AnalyzeImage)
		authenticated.GET("/nutrition/advice", nutritionHandler.GetAdvice)
		authenticated.GET("/nutrition/today-pfc", nutritionHandler.GetTodayPfc)
	}

	// Start server
	logger.Info("Starting server", "port", 8080)
	if err := r.Run(":8080"); err != nil {
		logger.Error("Failed to start server", "error", err.Error())
		panic(err)
	}
}
