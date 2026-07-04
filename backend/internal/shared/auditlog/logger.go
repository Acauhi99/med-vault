package auditlog

import (
	"context"
	"net/http"
	"strings"

	"github.com/Acauhi99/med-vault/internal/audit/application"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type Logger struct {
	logCmd *application.LogActionCommand
}

func NewLogger(logCmd *application.LogActionCommand) *Logger {
	return &Logger{logCmd: logCmd}
}

func (l *Logger) Log(ctx context.Context, r *http.Request, action, resourceType string, resourceID uuid.UUID, details map[string]any) {
	principal, ok := auth.PrincipalFromContext(ctx)

	var tenantID, userID uuid.UUID
	if ok {
		tenantID = principal.TenantID
		userID = principal.UserID
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	ua := r.UserAgent()

	go func() {
		_, _ = l.logCmd.Execute(context.Background(), application.LogActionInput{
			TenantID:     tenantID,
			UserID:       userID,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Details:      details,
			IPAddress:    ip,
			UserAgent:    ua,
		})
	}()
}
