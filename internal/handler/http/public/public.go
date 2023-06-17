package public

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/domain/websites"
	"github.com/qazaqpyn/webping/internal/handler/http/response"
	"github.com/qazaqpyn/webping/internal/service"
)

type HttpDelivery struct {
	resultsService service.Results
	websites       *websites.Websites
	auditService   service.Audit
}

func NewHttpDelivery(resultsService service.Results, websites *websites.Websites, auditService service.Audit) *HttpDelivery {
	return &HttpDelivery{
		resultsService: resultsService,
		websites:       websites,
		auditService:   auditService,
	}
}

type WebsiteRequest struct {
	Url string `json:"url"`
}

// @Summary GetRequestTime
// @Description	Get request time for specific website
// @Tags Public
// @Accept json
// @Produce	json
// @Param input body WebsiteRequest true "website url"
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/requestTime [post]
func (h *HttpDelivery) GetRequestTime(c *gin.Context) {
	r := WebsiteRequest{}
	if err := c.BindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	ok := h.websites.CheckWebExist(r.Url)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse("Website not pinged yet"))
		return
	}

	stat := h.resultsService.GetResult(r.Url)
	stat.ResponseTime = stat.ResponseTime / 1000000

	// store in audit
	if err := h.storeInAudit(c, audit.SPECIFIC_WEB, r.Url, stat.ResponseTime); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": stat,
	})
}

// @Summary GetMaxResponseTime
// @Description	Get max response time from all websites
// @Tags Public
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/maxResponseTime [post]
func (h *HttpDelivery) GetMaxResponseTime(c *gin.Context) {
	url, stat := h.resultsService.MaxResponseTime()
	stat.ResponseTime = stat.ResponseTime / 1000000

	// store in audit
	if err := h.storeInAudit(c, audit.MAX_RESPONSE_TIME, url, stat.ResponseTime); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"url":  url,
		"data": stat,
	})
}

// @Summary GetMinResponseTime
// @Description	Get min response time from all websites
// @Tags Public
// @Accept json
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400 {object} response.ResponseType
// @Failure	404	{object}	response.ResponseType
// @Failure	500	{object} response.ResponseType
// @Router /api/minResponseTime [post]
func (h *HttpDelivery) GetMinResponseTime(c *gin.Context) {
	url, stat := h.resultsService.MinResponseTime()
	stat.ResponseTime = stat.ResponseTime / 1000000

	// store in audit
	if err := h.storeInAudit(c, audit.MIN_RESPONSE_TIME, url, stat.ResponseTime); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]interface{}{
		"url":  url,
		"data": stat,
	})
}

func (h *HttpDelivery) storeInAudit(c *gin.Context, requestType int, url string, resTime time.Duration) error {
	return h.auditService.Create(c, requestType, url, resTime)
}
