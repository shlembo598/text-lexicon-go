package middleware

import (
	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/config"
)

// Middleware manager
type MiddlewareManager struct {
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string // список доступных доменов (CORS)
}

// Middleware manager constructor
func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, origins []string) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins}
}
