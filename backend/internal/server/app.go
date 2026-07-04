package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	auditapi "github.com/Acauhi99/med-vault/internal/audit"
	auditapp "github.com/Acauhi99/med-vault/internal/audit/application"
	pgxaudit "github.com/Acauhi99/med-vault/internal/audit/infrastructure/pgx"
	"github.com/Acauhi99/med-vault/internal/auth/application"
	"github.com/Acauhi99/med-vault/internal/auth/infrastructure/bcrypt"
	"github.com/Acauhi99/med-vault/internal/auth/infrastructure/jwt"
	pgxauth "github.com/Acauhi99/med-vault/internal/auth/infrastructure/pgx"
	clinicalapi "github.com/Acauhi99/med-vault/internal/clinical"
	pgxclinical "github.com/Acauhi99/med-vault/internal/clinical/infrastructure/pgx"
	"github.com/Acauhi99/med-vault/internal/generated"
	imagingapi "github.com/Acauhi99/med-vault/internal/imaging"
	imagingstorage "github.com/Acauhi99/med-vault/internal/imaging/infrastructure"
	pgximaging "github.com/Acauhi99/med-vault/internal/imaging/infrastructure/pgx"
	"github.com/Acauhi99/med-vault/internal/shared/auditlog"
	"github.com/Acauhi99/med-vault/internal/shared/config"
	"github.com/Acauhi99/med-vault/internal/shared/database"
	"github.com/Acauhi99/med-vault/internal/shared/httpx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config      config.Config
	DB          *pgxpool.Pool
	Logger      *slog.Logger
	Server      *http.Server
	AuditLogger *auditlog.Logger
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	db, err := database.OpenPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	hasher := bcrypt.NewHasher(cfg.BcryptCost)
	jwtGen := jwt.NewGenerator(cfg.JWTSecret)

	userRepo := pgxauth.NewUserRepository(db)
	tenantRepo := pgxauth.NewTenantRepository(db)

	registerCmd := application.NewRegisterCommand(userRepo, hasher)
	authenticateCmd := application.NewAuthenticateCommand(userRepo, tenantRepo, hasher, jwtGen, cfg.JWTTempTTL)
	selectTenantCmd := application.NewSelectTenantCommand(tenantRepo, jwtGen, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	refreshTokenCmd := application.NewRefreshTokenCommand(jwtGen, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	getCurrentUserQuery := application.NewGetCurrentUserQuery(userRepo, tenantRepo)
	addMemberCmd := application.NewAddMemberCommand(tenantRepo)
	removeMemberCmd := application.NewRemoveMemberCommand(tenantRepo)
	listMembersQuery := application.NewListMembersQuery(tenantRepo)
	reactivateTenantCmd := application.NewReactivateTenantCommand(tenantRepo)

	clinicalRepo := pgxclinical.NewCaseRepository(db)
	clinicalAPI := clinicalapi.NewAPI(clinicalRepo, tenantRepo)

	imagingRepo := pgximaging.NewImageRepository(db)
	storage := imagingstorage.NewStubStorage(cfg.S3Bucket, cfg.AWSRegion)
	imagingAPI := imagingapi.NewAPI(imagingRepo, clinicalRepo, storage)

	auditRepo := pgxaudit.NewAuditRepository(db)
	auditAPI := auditapi.NewAPI(auditRepo)
	logAuditCmd := auditapp.NewLogActionCommand(auditRepo)
	auditLog := auditlog.NewLogger(logAuditCmd)

	clinicalAPI.AuditLogger = auditLog
	imagingAPI.AuditLogger = auditLog

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
	})

	api := NewAPI(registerCmd, authenticateCmd, selectTenantCmd, refreshTokenCmd, getCurrentUserQuery, addMemberCmd, removeMemberCmd, listMembersQuery, reactivateTenantCmd, logger)
	api.ClinicalAPI = clinicalAPI
	api.ImagingAPI = imagingAPI
	api.AuditAPI = auditAPI
	api.auditLog = auditLog

	apiHandler := generated.HandlerWithOptions(
		api,
		generated.StdHTTPServerOptions{
			BaseURL: "/api/v1",
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			},
		},
	)
	mux.Handle("/api/v1/", apiHandler)

	authRateLimiter := httpx.NewRateLimiter(10, time.Minute)
	authedMux := http.NewServeMux()
	authedMux.Handle("/api/v1/auth/", authRateLimiter.Middleware(mux))
	authedMux.Handle("/", mux)
	handler := httpx.RequestIDMiddleware(cfg.RequestIDHeader)(authedMux)

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{Config: cfg, DB: db, Logger: logger, Server: httpServer, AuditLogger: auditLog}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
