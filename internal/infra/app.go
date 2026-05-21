package infra

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/viictormotorhead/financial-app/config"
	investmentUsecase "github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment"
	"github.com/viictormotorhead/financial-app/internal/infra/api/router/group"
	"github.com/viictormotorhead/financial-app/internal/infra/clients/db"
	"github.com/viictormotorhead/financial-app/internal/infra/logger"
	investmentGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/gorm"
)

func Run() {
	if err := config.Load(); err != nil {
		log.Panicf("error loading configuration: %v", err)
	}

	fx.New(
		fx.Provide(NewEchoServer),
		fx.Provide(NewEchoGroup),
		fx.Provide(logger.NewLogger),
		fx.Provide(db.NewPostgresConnection),
		fx.Provide(investmentGorm.NewInvestmentWriteRepository),
		fx.Provide(investmentUsecase.NewCreateInvestmentUseCase),
		fx.Provide(investment.NewInvestmentHandler),
		fx.Provide(group.NewHealthRoutes),
		fx.Provide(group.NewInvestmentRoutes),

		fx.Invoke(func(*echo.Echo) {}),
		fx.Invoke(runDatabaseMigrations),
		fx.Invoke(func(*group.HealthRoutes) {}),
		fx.Invoke(func(*group.InvestmentRoutes) {}),
	).Run()
}

func runDatabaseMigrations(gormDB *gorm.DB, logger *zap.Logger) error {
	if err := db.RunMigrations(gormDB); err != nil {
		return err
	}

	logger.Info("database migrations applied")
	return nil
}
