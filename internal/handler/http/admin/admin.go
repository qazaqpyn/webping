package admin

import (
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

// @Summary Login
// @Description	Login with Admin credentials
// @Tags Admin
// @Accept json
// @Produce	json
// @Param input body LoginRequest true "admin credentials"
// @Success	200	{object} map[string]string
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/login [post]
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

// @Summary GetWebList
// @Description	Get list of all API reqeusts from users to specific website
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/admin/webList [get]
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

// @Summary GetMinList
// @Description	Get list of all API reqeusts from users to minimum response time
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/admin/minList [get]
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

// @Summary GetMaxList
// @Description	Get list of all API reqeusts from users to maximum response time
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/admin/maxList [get]
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

// @Summary GetAllStatistics
// @Description	Get list of all API reqeusts from users (ID=1 : specific website, ID=2 : maximum response time, ID=3 : minimum response time)
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/admin/statisticAll [get]
func (h *HttpDelivery) GetAllStatistics(c *gin.Context) {
	stat, err := h.auditService.GetAll(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}
