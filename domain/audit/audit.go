package audit

import (
	"time"

	"github.com/qazaqpyn/webping/domain/base"
)

const (
	SPECIFIC_WEB = iota
	MAX_RESPONSE_TIME
	MIN_RESPONSE_TIME
)

type Audit struct {
	id           base.UUID
	requestType  int
	url          string
	ResponseTime time.Duration
}

func NewAudit(requestType int, url string, ResponseTime time.Duration) *Audit {
	auditId := base.NewUUID()

	return &Audit{
		id:           auditId,
		requestType:  requestType,
		url:          url,
		ResponseTime: ResponseTime,
	}
}

func (a *Audit) GetID() base.UUID {
	return a.id
}

func (a *Audit) GetRequestType() int {
	return a.requestType
}

func (a *Audit) GetURL() string {
	return a.url
}

func (a *Audit) GetResponseTime() time.Duration {
	return a.ResponseTime
}
