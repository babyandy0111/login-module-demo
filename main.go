package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	"github.com/shareed2k/goth_fiber"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	test := google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:8080/google/callback")
	// authCodeOptions []oauth2.AuthCodeOption

	goth.UseProviders(test)
	app.Get("/auth/:provider/login", goth_fiber.BeginAuthHandler)
	app.Get("/:provider/callback", func(ctx *fiber.Ctx) error {
		log.Println(ctx.Query("auid"))
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Fatal(err)
		}
		return ctx.JSON(fiber.Map{
			"email":   user.Email,
			"user_id": user.UserID,
		})
	})
	app.Get("/auth/:provider/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		return ctx.JSON(fiber.Map{
			"logout": "success",
		})
	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
