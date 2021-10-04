package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  viper.GetString("googleAuth.redirectURL"),
		ClientID:     viper.GetString("googleAuth.clientID"),
		ClientSecret: viper.GetString("googleAuth.clientSecret"),
		Scopes:       []string{viper.GetString("googleAuth.scopes")},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = viper.GetString("googleAuth.authStateString")
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

	token, err := googleOauthConfig.Exchange(context.Background(), c.Query("code"))

	if err != nil {
		fmt.Printf("could not get token %s\n", err.Error())
	}

	response, err := http.Get(viper.GetString("googleAuth.googleAuthUrl") + token.AccessToken)

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
