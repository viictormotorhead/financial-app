package infra

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/viictormotorhead/financial-app/internal/infra/config"
	investmentUsecase "github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	tagUsecase "github.com/viictormotorhead/financial-app/internal/application/tag/use_cases"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment"
	taghandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag"
	"github.com/viictormotorhead/financial-app/internal/infra/api/router/group"
	"github.com/viictormotorhead/financial-app/internal/infra/clients/db"
	"github.com/viictormotorhead/financial-app/internal/infra/logger"
	investmentGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/gorm"
	tagGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/gorm"
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
		fx.Provide(tagGorm.NewTagRepository),
		fx.Provide(investmentUsecase.NewCreateInvestmentUseCase),
		fx.Provide(investmentUsecase.NewListInvestmentsUseCase),
		fx.Provide(investmentUsecase.NewCreateMovementUseCase),
		fx.Provide(investmentUsecase.NewCreateValuationUseCase),
		fx.Provide(tagUsecase.NewCreateTagUseCase),
		fx.Provide(tagUsecase.NewListTagsUseCase),
		fx.Provide(investment.NewInvestmentHandler),
		fx.Provide(taghandler.NewTagHandler),
		fx.Provide(group.NewHealthRoutes),
		fx.Provide(group.NewInvestmentRoutes),
		fx.Provide(group.NewTagRoutes),

		fx.Invoke(func(*echo.Echo) {}),
		fx.Invoke(runDatabaseMigrations),
		fx.Invoke(func(*group.HealthRoutes) {}),
		fx.Invoke(func(*group.InvestmentRoutes) {}),
		fx.Invoke(func(*group.TagRoutes) {}),
	).Run()
}

func runDatabaseMigrations(gormDB *gorm.DB, logger *zap.Logger) error {
	if err := db.RunMigrations(gormDB); err != nil {
		return err
	}

	logger.Info("database migrations applied")
	return nil
}
