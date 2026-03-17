package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
	"github.com/kaecer68/lunar-zenith/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// 1. 初始化服務
	holidaySvc := service.NewHolidayService()
	// 嘗試載入 2024 樣本數據
	err := holidaySvc.LoadFromJSON("configs/holidays_2024_sample.json")
	if err != nil {
		log.Printf("Warning: Failed to load holiday data: %v", err)
	}

	aggregator := service.NewAggregator(holidaySvc)
	restHandler := service.NewRestHandler(aggregator)

	// 2. 設置 Gin (REST API)
	r := gin.Default()

	// 增加中間件處理繁體中文與安全性
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// 3. 註冊路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project": "Lunar-Zenith (算曆之巔)",
			"version": "v1.0.0",
			"status":  "running",
		})
	})

	restHandler.RegisterRoutes(r)

	// 4. 啟動 gRPC 服務器（在背景 goroutine）
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50051"
		}
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen gRPC: %v", err)
		}

		s := grpc.NewServer()
		grpcServer := service.NewGrpcServer(aggregator)
		lunarv1.RegisterLunarServiceServer(s, grpcServer)

		log.Printf("Lunar-Zenith gRPC starts on :%s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// 5. 啟動 REST HTTP 服務
	port := "8080"
	log.Printf("Lunar-Zenith REST API starts on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
