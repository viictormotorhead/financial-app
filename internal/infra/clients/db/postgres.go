package db

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/viictormotorhead/financial-app/config"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

func NewPostgresConnection(lc fx.Lifecycle, logger *zap.Logger) (*gorm.DB, error) {
	dbCfg := config.Config().DB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.Port,
		dbCfg.SSLMode,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database connection")
			return sqlDB.Close()
		},
	})

	return gormDB, nil
}

func RunMigrations(gormDB *gorm.DB) error {
	return gormDB.AutoMigrate(
		&entities.Tag{},
		&entities.Investment{},
		&entities.InvestmentHistory{},
	)
}
