package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaecer68/lunar-zenith/internal/service"
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

	// 2. 設置 Gin
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

	// 4. 啟動服務
	port := "8080"
	log.Printf("Lunar-Zenith API starts on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
