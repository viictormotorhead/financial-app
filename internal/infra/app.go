package infra

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	investmentUsecase "github.com/viictormotorhead/financial-app/internal/application/investment/use_cases"
	tagUsecase "github.com/viictormotorhead/financial-app/internal/application/tag/use_cases"
	userservices "github.com/viictormotorhead/financial-app/internal/application/user/services"
	userusecases "github.com/viictormotorhead/financial-app/internal/application/user/use_cases"
	authhandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler/auth"
	"github.com/viictormotorhead/financial-app/internal/infra/api/handler/investment"
	taghandler "github.com/viictormotorhead/financial-app/internal/infra/api/handler/tag"
	"github.com/viictormotorhead/financial-app/internal/infra/api/router/group"
	"github.com/viictormotorhead/financial-app/internal/infra/clients/db"
	"github.com/viictormotorhead/financial-app/internal/infra/config"
	"github.com/viictormotorhead/financial-app/internal/infra/logger"
	investmentGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/gorm"
	tagGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/gorm"
	userGorm "github.com/viictormotorhead/financial-app/internal/infra/repositories/user/gorm"
	"github.com/viictormotorhead/financial-app/internal/infra/security"
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
		fx.Provide(security.NewJWTTokenServiceFromConfig),
		fx.Provide(func(s *security.JWTTokenService) appauth.TokenIssuer { return s }),
		fx.Provide(func(s *security.JWTTokenService) appauth.TokenParser { return s }),
		fx.Provide(func() userservices.PasswordHasher {
			return security.NewBcryptHasher(config.Config().Auth.BcryptCost)
		}),
		fx.Provide(func() userservices.IDGenerator {
			return security.NewNanoidGenerator()
		}),
		fx.Provide(userGorm.NewUserRepository),
		fx.Provide(investmentGorm.NewInvestmentWriteRepository),
		fx.Provide(tagGorm.NewTagRepository),
		fx.Provide(userusecases.NewRegisterUserUseCase),
		fx.Provide(userusecases.NewLoginUserUseCase),
		fx.Provide(investmentUsecase.NewCreateInvestmentUseCase),
		fx.Provide(investmentUsecase.NewListInvestmentsUseCase),
		fx.Provide(investmentUsecase.NewGetInvestmentDetailUseCase),
		fx.Provide(investmentUsecase.NewCreateMovementUseCase),
		fx.Provide(investmentUsecase.NewCreateValuationUseCase),
		fx.Provide(tagUsecase.NewCreateTagUseCase),
		fx.Provide(tagUsecase.NewListTagsUseCase),
		fx.Provide(authhandler.NewAuthHandler),
		fx.Provide(investment.NewInvestmentHandler),
		fx.Provide(taghandler.NewTagHandler),
		fx.Provide(group.NewHealthRoutes),
		fx.Provide(group.NewAuthRoutes),
		fx.Provide(group.NewInvestmentRoutes),
		fx.Provide(group.NewTagRoutes),

		fx.Invoke(func(*echo.Echo) {}),
		fx.Invoke(func(*group.HealthRoutes) {}),
		fx.Invoke(func(*group.AuthRoutes) {}),
		fx.Invoke(func(*group.InvestmentRoutes) {}),
		fx.Invoke(func(*group.TagRoutes) {}),
	).Run()
}
