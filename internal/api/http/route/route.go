package route

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongnqse05461/ewallet/internal/api/http/handler"
	"github.com/truongnqse05461/ewallet/internal/api/http/middleware"
	"github.com/truongnqse05461/ewallet/internal/metrics"
)

type Route struct {
	metric             *metrics.Metric
	healthHandler      *handler.HealthHandler
	swaggerHandler     *handler.SwaggerHandler
	userHandler        *handler.UserHandler
	walletHandler      *handler.WalletHandler
	transactionHandler *handler.TransactionHandler
}

func NewRoute(
	metric *metrics.Metric,
	healthHandler *handler.HealthHandler,
	swaggerHandler *handler.SwaggerHandler,
	userHandler *handler.UserHandler,
	walletHandler *handler.WalletHandler,
	transactionHandler *handler.TransactionHandler,
) *Route {
	return &Route{
		metric:             metric,
		healthHandler:      healthHandler,
		swaggerHandler:     swaggerHandler,
		userHandler:        userHandler,
		walletHandler:      walletHandler,
		transactionHandler: transactionHandler,
	}
}

// nolint: funlen
func (r *Route) Index(server *gin.Engine, mw *middleware.Middleware) {
	server.GET("/readyz", r.healthHandler.Ready)
	server.GET("/livez", r.healthHandler.Live)
	server.GET("/metrics", gin.WrapH(r.metric.Handler()))
	server.GET("/swagger/*any", r.swaggerHandler.Swagger())
	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, errors.New("not found"))
	})

	root := server.Group("/")
	root.Use(mw.Logger(), mw.Trace(), mw.Metric(), mw.Access(), mw.ErrorHandle(), mw.Tx())
	{
		user := root.Group("/")
		{
			user.POST("/v1/users", r.userHandler.Create)
		}

		wallet := root.Group("/")
		{
			wallet.GET("/v1/wallets", r.walletHandler.List)
			wallet.POST("/v1/wallets", r.walletHandler.Create)
			wallet.POST("/v1/wallets/transfer", r.walletHandler.Transfer)
		}

		transaction := root.Group("/")
		{
			transaction.GET("/v1/transactions", r.transactionHandler.List)
		}
	}
}
