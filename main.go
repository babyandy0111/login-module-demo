package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"log"
	"os"
	"strconv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	callBaclUrl := "https://auth.codegenapps.com"
	app := fiber.New()
	googleConfig := google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), callBaclUrl+"/auth/google/callback")
	facebookConfig := facebook.New(os.Getenv("FACEBOOK_CLIENT_ID"), os.Getenv("FACEBOOK_CLIENT_SECRET"), callBaclUrl+"/auth/facebook/callback")
	goth.UseProviders(googleConfig, facebookConfig)

	// google login 因為需要丟參數
	// 因此需要轉止一次
	app.Get("/auth/:uid/:provider/login", func(ctx *fiber.Ctx) error {
		uid, err := Bin2hex(ctx.Params("uid"))
		if err != nil {
			return err
		}
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

	app.Get("/:provider/callback", func(ctx *fiber.Ctx) error {
		state, err := Hex2Bin(ctx.Query("state"))
		if err != nil {
			return err
		}
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

	app.Get("/:provider/logout", func(ctx *fiber.Ctx) error {
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
func Hex2Bin(hex string) (string, error) {
	ui, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%016b", ui), nil
}

func Bin2hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 16), nil
}
