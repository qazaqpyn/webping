package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/webping/internal/handler/http/admin"
	"github.com/qazaqpyn/webping/internal/handler/http/middleware"
	"github.com/qazaqpyn/webping/internal/handler/http/public"
	"github.com/qazaqpyn/webping/internal/service"
)

type Handler struct {
	Admin  *admin.HttpDelivery
	Public *public.HttpDelivery
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Admin:  admin.NewHttpDelivery(services.Audit, services.Auth),
		Public: public.NewHttpDelivery(services.Results),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/login", h.Admin.Login)

		authentication := api.Group("/admin")
		authentication.Use(middleware.AdminIdentityMiddleware(h.Admin.AuthService))
		{
			authentication.GET("/statisticAll", h.Admin.GetAllStatistics)
			authentication.GET("/webStatistic", h.Admin.GetWebStatistics)
			authentication.GET("/minStatistic", h.Admin.GetMinStatistics)
			authentication.GET("/maxStatistic", h.Admin.GetMaxStatistics)
		}

		api.POST("/requestTime", h.Public.GetRequestTime)
		api.POST("/maxResponseTime", h.Public.GetMaxResponseTime)
		api.POST("/minResponseTime", h.Admin.GetMinStatistics)
	}
	return router
}
