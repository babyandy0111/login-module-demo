package main

import (
	"fmt"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/dghubble/sessions"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	"log"
	"net/http"
)

const (
	sessionName     = "example-google-app"
	sessionSecret   = "I Am CodeGenApps sessionSecret"
	sessionUserKey  = "googleID"
	sessionUsername = "googleName"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

type GoogleLoginConfig struct {
	ClientID        string
	ClientSecret    string
	UserRedirectURL string
}

func initGoogleLogin(oauth2Config *oauth2.Config, stateConfig gologin.CookieConfig) http.Handler {
	return google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil))
}

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = googleUser.Id
		session.Values[sessionUsername] = googleUser.Name
		session.Save(w)
		fmt.Println(googleUser.Email)
		fmt.Println(googleUser.Id)
		// http.Redirect(w, req, "https://s3.amazonaws.com/assets.codegenapps.com/index.html", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func googleLogin(oauth2Config *oauth2.Config, stateConfig gologin.CookieConfig) http.Handler {
	return google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))
}

func main() {
	// 1. 這個值
	config := &GoogleLoginConfig{
		ClientID:        "219033873816-tvtqlles09pelgknl2uepbjs0bm48kmv.apps.googleusercontent.com",
		ClientSecret:    "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",
		UserRedirectURL: "",
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://localhost:8080/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig

	app := fiber.New()

	app.Get("/:login", adaptor.HTTPHandler(googleLogin(oauth2Config, stateConfig)))
	app.Get("/:login/callback", adaptor.HTTPHandler(initGoogleLogin(oauth2Config, stateConfig)))

	log.Fatal(app.Listen(":8080"))
}
