package security

import "github.com/viictormotorhead/financial-app/internal/infra/config"

func NewJWTTokenServiceFromConfig() *JWTTokenService {
	authCfg := config.Config().Auth
	return NewJWTTokenService(authCfg.JWTSecret, authCfg.JWTExpiry)
}
