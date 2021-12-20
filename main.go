package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/line"
	"github.com/shareed2k/goth_fiber"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	callBaclUrl := "http://localhost:8080"
	app := fiber.New()
	googleConfig := google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), callBaclUrl+"/auth/google/callback")
	facebookConfig := facebook.New(os.Getenv("FACEBOOK_CLIENT_ID"), os.Getenv("FACEBOOK_CLIENT_SECRET"), callBaclUrl+"/auth/facebook/callback")
	appleConfig := apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), callBaclUrl+"/auth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail)
	lineConfig := line.New(os.Getenv("LINE_KEY"), os.Getenv("LINE_SECRET"), callBaclUrl+"/auth/line/callback", "profile", "openid", "email")
	goth.UseProviders(googleConfig, facebookConfig, lineConfig, appleConfig)

	// google login 因為需要丟參數
	// 因此需要轉止一次
	app.Get("/auth/:uid/:provider/login", func(ctx *fiber.Ctx) error {
		uid := ctx.Params("uid")
		log.Println(uid)
		newURL := fmt.Sprintf("/auth/%s?state=%s", ctx.Params("provider"), uid)
		log.Println(newURL)
		return ctx.Redirect(newURL, fiber.StatusTemporaryRedirect)
	})

	//
	app.Get("/auth/:provider", func(ctx *fiber.Ctx) error {
		url, err := goth_fiber.GetAuthURL(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		log.Println(url)
		return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
	})

	app.Get("/auth/:provider/callback", func(ctx *fiber.Ctx) error {
		state := ctx.Query("state")
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			return err
		}

		log.Println("state", state)
		clitneURL := "https://s3.amazonaws.com/assets.codegenapps.com/index.html"
		token := "QQ"
		newURL := fmt.Sprintf("%s?token=%s&gid=%s&email=%s", clitneURL, token, user.UserID, user.Email)

		return ctx.Redirect(newURL, fiber.StatusTemporaryRedirect)
	})

	app.Get("/auth/:provider/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			return err
		}
		clitneURL := "https://s3.amazonaws.com/assets.codegenapps.com/index.html"
		return ctx.Redirect(clitneURL, fiber.StatusTemporaryRedirect)
	})

	if err := app.Listen("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
