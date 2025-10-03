package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getDashboardSummary 대시보드 요약 정보 조회
func (app *Application) getDashboardSummary(c *gin.Context) {
	summary := generateMockDashboardSummary()

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

// getTraders 트레이더 목록 조회
func (app *Application) getTraders(c *gin.Context) {
	var req TradersRequest

	// 쿼리 파라미터 파싱
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// 기본값 설정
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 50
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	// Mock 데이터 생성
	response := getMockTradersWithPagination(req.Page, req.PerPage, req.Search)

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// getTraderDetail 특정 트레이더 상세 정보 조회
func (app *Application) getTraderDetail(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address parameter is required",
		})
		return
	}

	// Mock 데이터에서 해당 주소의 트레이더 찾기
	allTraders := generateMockTraders(1000)
	var foundTrader *Trader

	for _, trader := range allTraders {
		if trader.Address == address {
			foundTrader = &trader
			break
		}
	}

	if foundTrader == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Trader not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": foundTrader,
	})
}

// getTraderStats 트레이더 통계 정보 조회
func (app *Application) getTraderStats(c *gin.Context) {
	// 전체 트레이더 데이터 생성
	allTraders := generateMockTraders(1000)

	// 통계 계산
	var totalEquity float64
	var profitableCount int
	var totalDailyPnL float64

	for _, trader := range allTraders {
		totalEquity += trader.PerpEquity
		if trader.DailyPnL.Amount > 0 {
			profitableCount++
		}
		totalDailyPnL += trader.DailyPnL.Amount
	}

	avgDailyPnL := totalDailyPnL / float64(len(allTraders))
	profitablePercentage := float64(profitableCount) / float64(len(allTraders)) * 100

	stats := gin.H{
		"total_traders":         len(allTraders),
		"total_equity":          totalEquity,
		"profitable_count":      profitableCount,
		"profitable_percentage": profitablePercentage,
		"avg_daily_pnl":         avgDailyPnL,
		"total_daily_pnl":       totalDailyPnL,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// getAssetStats 자산별 통계 조회
func (app *Application) getAssetStats(c *gin.Context) {
	allTraders := generateMockTraders(1000)

	// 자산별 집계
	assetStats := make(map[string]map[string]interface{})

	for _, trader := range allTraders {
		if trader.MainPosition != nil {
			asset := trader.MainPosition.Asset

			if _, exists := assetStats[asset]; !exists {
				assetStats[asset] = map[string]interface{}{
					"asset":         asset,
					"total_traders": 0,
					"total_oi":      0.0,
					"avg_position":  0.0,
					"long_count":    0,
					"short_count":   0,
				}
			}

			stats := assetStats[asset]
			stats["total_traders"] = stats["total_traders"].(int) + 1
			stats["total_oi"] = stats["total_oi"].(float64) + trader.MainPosition.Amount

			if trader.MainPosition.Type == "LONG" {
				stats["long_count"] = stats["long_count"].(int) + 1
			} else {
				stats["short_count"] = stats["short_count"].(int) + 1
			}
		}
	}

	// 평균 포지션 계산
	for _, stats := range assetStats {
		if stats["total_traders"].(int) > 0 {
			stats["avg_position"] = stats["total_oi"].(float64) / float64(stats["total_traders"].(int))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": assetStats,
	})
}
