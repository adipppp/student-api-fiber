package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"studentapifiber/constants"
	"studentapifiber/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gopkg.in/cas.v2"
)

var casClient *cas.Client

func NewAuthMiddlewareHandler() (fiber.Handler, error) {
	authHandler, err := newAuthHandler()
	if err != nil {
		return nil, fmt.Errorf("error creating auth handler: %v", err)
	}
	authFiberHandler := adaptor.HTTPHandler(authHandler)

	return func(c *fiber.Ctx) error {
		authFiberHandler(c)
		if c.Response().StatusCode() == http.StatusUnauthorized {
			return nil
		}
		return c.Next()
	}, nil
}

func getCasClient() (*cas.Client, error) {
	if casClient != nil {
		return casClient, nil
	}

	client, err := newCasClient()
	if err != nil {
		return nil, fmt.Errorf("error creating CAS client: %v", err)
	}

	return client, nil
}

func newAuthHandler() (http.Handler, error) {
	client, err := getCasClient()
	if err != nil {
		return nil, fmt.Errorf("error getting CAS client: %v", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		if !cas.IsAuthenticated(r) {
			errorResponse := dto.ErrorResponse{
				Code:    constants.Unauthorized,
				Message: "Unauthorized, please login via SSO",
			}
			w.WriteHeader(http.StatusUnauthorized)
			enc.Encode(errorResponse)
		}
	}

	return client.HandleFunc(handler), nil
}

func newCasClient() (*cas.Client, error) {
	ssoUrlString, ok := os.LookupEnv("SSO_UI_URL")
	if !ok {
		return nil, fmt.Errorf("SSO_UI_URL environment variable not set")
	}

	ssoUrl, err := url.Parse(ssoUrlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing SSO_UI_URL: %v", err)
	}

	client := cas.NewClient(&cas.Options{URL: ssoUrl})
	casClient = client

	return client, nil
}
