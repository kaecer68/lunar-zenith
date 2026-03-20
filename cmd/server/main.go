package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
	"github.com/kaecer68/lunar-zenith/internal/service"
	"github.com/kaecer68/lunar-zenith/internal/webui"
	"google.golang.org/grpc"
)

func getEnvWithFallback(primary, fallback, defaultValue string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		return v
	}
	return defaultValue
}

func main() {
	// 1. 初始化服務
	// 台灣假期 - 使用完整假期數據
	holidaySvc := service.NewHolidayService()
	if err := holidaySvc.LoadFromJSON("configs/holidays_tw_2024_2026.json"); err != nil {
		log.Printf("Warning: Failed to load Taiwan holiday data: %v", err)
	}

	// 大陸假期
	chinaHolidaySvc := service.NewHolidayService()
	if err := chinaHolidaySvc.LoadFromJSON("configs/holidays_cn_2024_2026.json"); err != nil {
		log.Printf("Warning: Failed to load China holiday data: %v", err)
	}

	aggregator := service.NewAggregator(holidaySvc, chinaHolidaySvc)
	restHandler := service.NewRestHandler(aggregator)

	// 2. 設置 Gin (REST API)
	r := gin.Default()

	// 3. 靜態前端 (嵌入 binary)
	r.StaticFS("/ui", webui.FS())
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/")
	})

	// 4. 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project": "Lunar-Zenith (算曆之巔)",
			"version": "v1.4.0",
			"status":  "running",
		})
	})

	restHandler.RegisterRoutes(r)

	// 4. 啟動 gRPC 服務器（在背景 goroutine）
	go func() {
		grpcPort := getEnvWithFallback("LUNAR_GRPC_PORT", "GRPC_PORT", "50051")
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
	port := getEnvWithFallback("LUNAR_REST_PORT", "REST_PORT", "8080")
	log.Printf("Lunar-Zenith REST API starts on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
