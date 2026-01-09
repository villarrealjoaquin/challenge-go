package handlers

import (
	"net/http"

	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

type GetMetricsRequest struct {
	Author string `form:"author" binding:"required"`
}

func NewGetMetrics(metricsService services.MetricsService) GetMetrics {
	return GetMetrics{metricsService}
}

type GetMetrics struct {
	metricsService services.MetricsService
}

func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query GetMetricsRequest
		if err := ctx.ShouldBindQuery(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid query parameters",
			})
			return
		}

		metrics, err := h.metricsService.GetMetrics(ctx.Request.Context(), query.Author)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to get metrics",
			})
			return
		}

		ctx.JSON(http.StatusOK, metrics)
	}
}
