package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/qazaqpyn/webping/docs"
	"github.com/qazaqpyn/webping/domain/websites"
	"github.com/qazaqpyn/webping/internal/handler/http/admin"
	"github.com/qazaqpyn/webping/internal/handler/http/middleware"
	"github.com/qazaqpyn/webping/internal/handler/http/public"
	"github.com/qazaqpyn/webping/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Admin  *admin.HttpDelivery
	Public *public.HttpDelivery
}

func NewHandler(services *service.Service, websites *websites.Websites) *Handler {
	return &Handler{
		Admin:  admin.NewHttpDelivery(services.Audit, services.Auth),
		Public: public.NewHttpDelivery(services.Results, websites, services.Audit),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/login", h.Admin.Login)
		api.POST("/requestTime", h.Public.GetRequestTime)
		api.POST("/maxResponseTime", h.Public.GetMaxResponseTime)
		api.POST("/minResponseTime", h.Public.GetMinResponseTime)

		authentication := api.Group("/admin")
		authentication.Use(middleware.AdminIdentityMiddleware(h.Admin.AuthService))
		{
			authentication.GET("/statisticAll", h.Admin.GetAllStatistics)
			authentication.GET("/webList", h.Admin.GetWebList)
			authentication.GET("/minList", h.Admin.GetMinList)
			authentication.GET("/maxList", h.Admin.GetMaxList)
		}
	}
	return router
}
