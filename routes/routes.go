package routes


import (
	"github.com/lokesh2201013/email-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/verify-email-identity", controllers.VerifyEmailID)
	app.Get("/list-verified-identities",  controllers.ListVerifiedIdentities)
	app.Get("/list-unverified-identities",  controllers.ListunVerifiedIdentities)
	app.Delete("/delete-identity/:email",  controllers.DeleteIdentity)
	app.Post("/send-email",  controllers.SendEmail)

	app.Post("/create-template", controllers.CreateTemplate)
	//app.Post("/send-custom-email",  controllers.SendCustomEmail)
}