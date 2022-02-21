package api

import (
	"accounts/application"
	"accounts/pkg/api/handlers"
	"accounts/pkg/api/middlewares"
	"accounts/pkg/cache"
	"accounts/pkg/db"
	"accounts/pkg/encrypter"
	"accounts/pkg/queue"
	"database/sql"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/streadway/amqp"
)

func SetupServer(router *echo.Echo, database *sql.DB, rabbitmq *amqp.Channel, client *redis.Client) {

	// init adapters
	accountRepo := db.NewPgAccountRepository(database)
	redisRepo := cache.NewRedisCacheRepository(client)
	jwtEncrypter := encrypter.NewEncrypter(redisRepo)
	rabbitMQ := queue.NewRabbitMQ(rabbitmq)

	// init usecases
	createAccount := application.NewCreateAccount(accountRepo)
	authAccount := application.NewAuthAccount(accountRepo, jwtEncrypter, redisRepo)
	refreshAuth := application.NewRefreshAuth(jwtEncrypter, accountRepo, redisRepo)
	deleteAccount := application.NewDeleteUsecase(accountRepo)
	findAccount := application.NewFindAccount(accountRepo)
	updateAccount := application.NewUpdateAccount(accountRepo)

	// init handlers
	createAccountHandler := handlers.NewCreateAccountHandler(*createAccount, rabbitMQ, jwtEncrypter)
	authHandler := handlers.NewAuthHandler(authAccount)
	refreshAuthHandler := handlers.NewRefreshAuthHandler(refreshAuth)
	deleteAccountHandler := handlers.NewDeleteAccountHandler(deleteAccount, rabbitMQ)
	getAccountHandler := handlers.NewFindAccountHandler(findAccount)
	updateAccountHandler := handlers.NewUpdateAccountHandler(updateAccount, rabbitMQ)

	// init middlewares
	authMiddleware := middlewares.NewAuthMiddleware(jwtEncrypter)

	// setup middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(middleware.Secure())
	router.Use(middleware.BodyLimit("2M"))
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// setup routes
	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "Ok"})
	})
	router.GET("/error", func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "Err"})
	})

	r := router.Group("/accounts")

	r.POST("", createAccountHandler.Handle)
	r.POST("/auth", authHandler.Handle)
	r.POST("/auth/refresh", refreshAuthHandler.Handle)
	r.GET("/:id", getAccountHandler.Handle, authMiddleware.Auth)
	r.DELETE("/:id", deleteAccountHandler.Handle, authMiddleware.Auth)
	r.PUT("/:id", updateAccountHandler.Handle, authMiddleware.Auth)
}
