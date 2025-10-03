package apiserver

import (
	"net/http"

	"github.com/cresendoo/decidash-backend/internal/application/api-server/middleware"
	"github.com/gin-gonic/gin"
)

func (app *Application) setRouter() http.Handler {
	handler := gin.New()

	handler.Use(middleware.CORS())

	handler.GET("health_check", middleware.HealthCheck())

	api := handler.Group("/api")
	apiV1 := api.Group("/v1")

	// set middleware
	apiV1.Use(
		middleware.SetRequestID,
		middleware.SetRequsetLogger(app.logger),
		middleware.GinRecovery(),
		middleware.RequestLog,
		middleware.ResponseHandler,
	)
	apiV1.GET("/test_error", middleware.TestError())

	traders := apiV1.Group("/traders")
	{
		traders.GET("/dashboard", app.getDashboardSummary)
		traders.GET("", app.getTraders)
		traders.GET("/:address", app.getTraderDetail)
		traders.GET("/stats", app.getTraderStats)
		traders.GET("/assets/stats", app.getAssetStats)
	}

	transactions := apiV1.Group("/transactions")
	{
		transactions.POST("", app.postFeePayer)
	}
	return handler
}
