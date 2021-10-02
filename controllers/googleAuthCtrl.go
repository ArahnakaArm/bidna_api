package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3334/apix/v1/google_auth/callback",
		ClientID:     "781374464096-1gdncta7qoh9atlpf07sb4gfonq04lci.apps.googleusercontent.com",
		ClientSecret: "GaXVKMIHTuT-dNTkEBMym-36",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "pseudo-random"
)

func HandleGoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func HandleGoogleCallBack(c *fiber.Ctx) error {
	if c.Query("state") != oauthStateString {
		fmt.Printf("state is not valid")
		return c.Redirect("/", fiber.StatusTemporaryRedirect)
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, c.Query("code"))

	if err != nil {
		fmt.Printf("could not get token %s\n", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		fmt.Println("could not get request")
		return c.Redirect("/", fiber.StatusTemporaryRedirect)
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("could not parse ")
		return c.Redirect("/", fiber.StatusTemporaryRedirect)
	}

	return c.Send(content)
}
