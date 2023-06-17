package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/internal/handler/http/response"
	"github.com/qazaqpyn/webping/internal/service"
)

type HttpDelivery struct {
	auditService service.Audit
	AuthService  service.Auth
}

func NewHttpDelivery(service service.Audit, authService service.Auth) *HttpDelivery {
	return &HttpDelivery{auditService: service, AuthService: authService}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HttpDelivery) Login(c *gin.Context) {
	r := LoginRequest{}
	if err := c.BindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	token, err := h.AuthService.Login(c, r.Email, r.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *HttpDelivery) GetWebList(c *gin.Context) {
	stat, err := h.auditService.GetByRequestType(c, audit.SPECIFIC_WEB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}

func (h *HttpDelivery) GetMinList(c *gin.Context) {
	stat, err := h.auditService.GetByRequestType(c, audit.MIN_RESPONSE_TIME)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}

func (h *HttpDelivery) GetMaxList(c *gin.Context) {
	stat, err := h.auditService.GetByRequestType(c, audit.MAX_RESPONSE_TIME)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}

func (h *HttpDelivery) GetAllStatistics(c *gin.Context) {
	stat, err := h.auditService.GetAll(c)
	log.Println(stat, err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}
