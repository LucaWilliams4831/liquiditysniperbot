package routes

import (
	"github.com/LucaWilliams4831/uniswap-pancakeswap-tradingbot/liquiditysniperbot/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// User routes
	app.Post("/api/v1/user/login", controllers.Login)
	app.Post("/api/v1/user/logout", controllers.Logout)


	//Account routes
	app.Get("api/v1/account/get_all", controllers.GetAccounts)
	app.Get("api/v1/account/filter/:keyword", controllers.FilterAccounts)
	app.Post("api/v1/account/update/:id", controllers.UpdateAccount)
	app.Post("api/v1/account/sendfee/:id", controllers.SendFee)
	app.Post("api/v1/account/add/:address", controllers.AddAccount)

}
