package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RestHandler 處理 HTTP 請求
type RestHandler struct {
	Aggregator *Aggregator
}

// NewRestHandler 創建 REST 處理器
func NewRestHandler(agg *Aggregator) *RestHandler {
	return &RestHandler{Aggregator: agg}
}

// RegisterRoutes 註冊路由
func (h *RestHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/v1/calendar", h.GetCalendar)
}

// GetCalendar 獲取曆法數據
// Query: ?date=2024-03-14
func (h *RestHandler) GetCalendar(c *gin.Context) {
	dateStr := c.Query("date")
	var t time.Time
	var err error

	if dateStr == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
	}

	res := h.Aggregator.GetCalendar(t)
	c.JSON(http.StatusOK, res)
}
