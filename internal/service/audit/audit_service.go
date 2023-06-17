package audit

import (
	"context"
	"time"

	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/internal/repository"
)

type AuditService struct {
	auditRepo repository.Audit
}

func NewAuditService(repo repository.Audit) *AuditService {
	return &AuditService{auditRepo: repo}
}

func (s *AuditService) Create(ctx context.Context, requestType int, url string, ResponseTime time.Duration) error {
	newUser := audit.NewAudit(requestType, url, ResponseTime)

	err := s.auditRepo.Create(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuditService) GetAll(ctx context.Context) ([]*audit.MongoAuditGroup, error) {
	return s.auditRepo.Find(ctx)
}

func (s *AuditService) GetByRequestType(ctx context.Context, requestType int) ([]*audit.MongoAuditResp, error) {
	return s.auditRepo.FindByRequestType(ctx, requestType)
}
