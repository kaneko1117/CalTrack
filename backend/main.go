package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"caltrack/config"
	"caltrack/handler/user"
	gormPersistence "caltrack/infrastructure/persistence/gorm"
	"caltrack/pkg/logger"
	"caltrack/usecase"
)

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

	// DI
	userRepo := gormPersistence.NewGormUserRepository(config.DB)
	txManager := gormPersistence.NewGormTransactionManager(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, txManager)
	userHandler := user.NewUserHandler(userUsecase)

	// Setup router
	r := gin.Default()

	// ロガーミドルウェア追加
	r.Use(logger.Middleware())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "CalTrack API is running"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.Register)
	}

	// Start server
	logger.Info("Starting server", "port", 8080)
	if err := r.Run(":8080"); err != nil {
		logger.Error("Failed to start server", "error", err.Error())
		panic(err)
	}
}
