package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "hw4/docs"
	"hw4/internal/app/pkg/auth"
	core "hw4/internal/app/pkg/core"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/dish"
	pgdish "hw4/internal/app/pkg/dish/postgresql"
	"hw4/internal/app/pkg/order"
	pgorder "hw4/internal/app/pkg/order/postgresql"
	"hw4/internal/app/pkg/order_dish"
	pgorderdish "hw4/internal/app/pkg/order_dish/postgresql"
	"hw4/internal/app/pkg/order_processing"
	"hw4/internal/app/pkg/session"
	pgsession "hw4/internal/app/pkg/session/postgresql"
	"hw4/internal/app/pkg/user"
	pguser "hw4/internal/app/pkg/user/postgresql"
	"net/http"
)

func HelloWorld(c *gin.Context) {
	fmt.Println("Hello, World!")
	c.String(http.StatusOK, "hi")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//	@title			Swagger API
//	@version		1.0
//	@description	UI for microservers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		fmt.Println("db connection error")
		return
	}
	defer database.GetPool(ctx).Close()

	dishRepo := pgdish.NewDish(database)
	orderRepo := pgorder.NewOrder(database)
	orderDishRepo := pgorderdish.NewOrderDish(database)
	sessionRepo := pgsession.NewSession(database)
	userRepo := pguser.NewUser(database)

	core := core.NewService(
		dish.NewService(dishRepo),
		order.NewService(orderRepo),
		order_dish.NewService(orderDishRepo),
		session.NewService(sessionRepo),
		user.NewService(userRepo))

	authService := auth.NewService(core)
	orderProcessingService := order_processing.NewService(core)

	router := gin.New()

	router.Use(CORSMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/get-user-info", authService.GetUserInfo)
		auth.POST("/register-new-user", authService.RegisterNewUser)
		auth.POST("/authorization", authService.Authorization)
	}

	//router.POST("/get-user-info", authService.GetUserInfo)
	//router.POST("/register-new-user", authService.RegisterNewUser)
	//router.POST("/authorization", authService.Authorization)

	orderProcessing := router.Group("/order-processing")
	{
		orderProcessing.GET("/get-menu", orderProcessingService.GetMenu)
		orderProcessing.POST("/get-order", orderProcessingService.GetOrder)
		orderProcessing.POST("/get-dish", orderProcessingService.GetDish)
		orderProcessing.POST("/create-order", orderProcessingService.CreateOrder)
		orderProcessing.POST("/add-dish", orderProcessingService.AddDish)
		orderProcessing.POST("/delete-dish", orderProcessingService.DeleteDish)
		orderProcessing.PUT("/update-order", orderProcessingService.UpdateOrder)
		orderProcessing.PUT("/update-dish", orderProcessingService.UpdateDish)
	}

	router.Run(":8080")
}
