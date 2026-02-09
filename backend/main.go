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

	// Connect to database
	config.ConnectDatabase()

	// Geminiクライアント初期化
	config.InitGemini()
	defer config.CloseGemini()

	// DI - Repository
	userRepo := gormPersistence.NewGormUserRepository(config.DB)
	sessionRepo := gormPersistence.NewGormSessionRepository(config.DB)
	recordRepo := gormPersistence.NewGormRecordRepository(config.DB)
	recordPfcRepo := gormPersistence.NewGormRecordPfcRepository(config.DB)
	adviceCacheRepo := gormPersistence.NewGormAdviceCacheRepository(config.DB)
	txManager := gormPersistence.NewGormTransactionManager(config.DB)

	// DI - Service
	imageAnalyzer := infraService.NewGeminiImageAnalyzer(config.GeminiClient)
	pfcAnalyzer := infraService.NewGeminiPfcAnalyzer(config.GeminiClient)
	pfcEstimator := infraService.NewGeminiPfcEstimator(config.GeminiClient)

	// DI - Config
	aiConfig := config.DefaultAIConfig{}

	// DI - Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, txManager)
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	recordUsecase := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, aiConfig)
	analyzeUsecase := usecase.NewAnalyzeUsecase(imageAnalyzer, aiConfig)
	nutritionUsecase := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, pfcAnalyzer, aiConfig)

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
