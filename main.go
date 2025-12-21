package main

import (
	"log"

	"github.com/ijufumi/practice-202512/app/infrastructure/database/gateway"

	"github.com/ijufumi/practice-202512/app/config"
	"github.com/ijufumi/practice-202512/app/infrastructure/database"
	"github.com/ijufumi/practice-202512/app/presentation"
	"github.com/ijufumi/practice-202512/app/presentation/handler"
	"github.com/ijufumi/practice-202512/app/usecase"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

	// データベース接続
	db, err := database.NewConnection(cfg)
	log.Println("test")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 依存性の注入
	invoiceRepository := gateway.NewInvoiceRepository()
	userRepository := gateway.NewUserRepository()
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepository, userRepository)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUsecase)

	authUsecase := usecase.NewAuthUsecase(userRepository, cfg)
	authHandler := handler.NewAuthHandler(authUsecase)

	// ルーター設定
	router := presentation.NewRouter(db, cfg, invoiceHandler, authHandler)
	defer func() {
		_ = router.Close()
	}()

	// サーバー起動
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
