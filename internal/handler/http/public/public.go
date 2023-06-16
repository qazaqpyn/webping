package public

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/webping/domain/websites"
	"github.com/qazaqpyn/webping/internal/handler/http/response"
	"github.com/qazaqpyn/webping/internal/service"
)

type HttpDelivery struct {
	resultsService service.Results
	websites       *websites.Websites
}

func NewHttpDelivery(resultsService service.Results, websites *websites.Websites) *HttpDelivery {
	return &HttpDelivery{
		resultsService: resultsService,
		websites:       websites,
	}
}

type WebsiteRequest struct {
	Url string `json:"url"`
}

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

	data, err := json.Marshal(map[string]interface{}{
		"data": stat,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}

func (h *HttpDelivery) GetMaxResponseTime(c *gin.Context) {
	url, stat := h.resultsService.MaxResponseTime()

	data, err := json.Marshal(map[string]interface{}{
		"url":  url,
		"data": stat,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}

func (h *HttpDelivery) GetMinResponseTime(c *gin.Context) {
	url, stat := h.resultsService.MinResponseTime()

	data, err := json.Marshal(map[string]interface{}{
		"url":  url,
		"data": stat,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}
