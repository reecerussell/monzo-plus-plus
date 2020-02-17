package monzo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/model"
)

// Environment variables.
var (
	MonzoClientID      = os.Getenv("MONZO_CLIENT_ID")
	MonzoClientSecret  = os.Getenv("MONZO_CLIENT_SECRET")
	OAuthCallBackURL   = os.Getenv("OAUTH_CALLBACK_URL")
	SuccessCallbackURL = os.Getenv("SUCCESS_CALLBACK_URL")
)

// Monzo status code errors
var (
	ErrBadRequest          = errors.New("your request has missing arguments or is malformed")
	ErrUnauthorized        = errors.New("your request is not authenticated")
	ErrForbidden           = errors.New("your request is authenticated but has insufficent permissions")
	ErrMethodNotAllowed    = errors.New("you are using an incorrect HTTP verb. double check whether it should be POST/GET/DELETE/etc")
	ErrPageNotFound        = errors.New("the endpoint requested does not exist")
	ErrNotAcceptable       = errors.New("your application does not accept the content format returned according to the accpet headers sent in the request")
	ErrTooManyRequests     = errors.New("your application is exceeding its rate limit. back off, buddy :p")
	ErrInternalServerError = errors.New("something is wrong on our end. whoopsie")
	ErrGatewayTimeout      = errors.New("something has timed out on our end. whoopsie")
)

// Monzo grant types.
const (
	GrantTypeAuthCode     = "authorization_code"
	GrantTypeRefreshToken = "refresh_token"
)

var (
	defaultClient Client
)

func init() {
	defaultClient = NewClient()
}

// Client is a high level service interface used for making requests to monzo.
type Client interface {
	IsAuthenticated(u *model.User) bool
	Login(w http.ResponseWriter, r *http.Request, u *model.User)
	RequestAccessToken(code string) (*AccessToken, error)
	RefreshAccessToken(refreshToken string) (*AccessToken, error)
	Logout(u *model.User) error
}

func mustUseDefaultClient() Client {
	if defaultClient == nil {
		panic("default client must be set before use")
	}

	return defaultClient
}

func IsAuthenticated(u *model.User) bool {
	return mustUseDefaultClient().IsAuthenticated(u)
}

func Login(w http.ResponseWriter, r *http.Request, u *model.User) {
	mustUseDefaultClient().Login(w, r, u)
}

func Logout(u *model.User) error {
	return nil
}

func RequestAccessToken(code string) (*AccessToken, error) {
	return mustUseDefaultClient().RequestAccessToken(code)
}

func RefreshAccessToken(refreshToken string) (*AccessToken, error) {
	return mustUseDefaultClient().RefreshAccessToken(refreshToken)
}

type client struct {
	h *http.Client
}

// NewClient returns a new instance of client.
func NewClient() Client {
	return &client{
		h: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *client) IsAuthenticated(u *model.User) bool {
	log.Printf("Attempting to validate user's authenticity.")
	ut := u.GetToken()
	ac := ut.GetAccessToken()
	if ac == "" {
		log.Printf("\tAccess token is empty.\n")
		return false
	}

	req, _ := http.NewRequest(http.MethodGet, "https://api.monzo.com/ping/whoami", nil)
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", ut.GetTokenType(), ac))

	log.Printf("\tMaking whoami request...\n")
	resp, _ := c.h.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e Error
		_ = json.NewDecoder(resp.Body).Decode(&e)
		log.Printf("\t\tFailed: %s\n", e.Message)
		return false
	}

	log.Printf("\tMade request successfully.")

	var ad AuthenticationData
	_ = json.NewDecoder(resp.Body).Decode(&ad)

	return ad.Authenticated
}

func (c *client) Login(w http.ResponseWriter, r *http.Request, u *model.User) {
	url := fmt.Sprintf("https://auth.monzo.com/?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		MonzoClientID,
		OAuthCallBackURL,
		u.GetStateToken())

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *client) RequestAccessToken(code string) (*AccessToken, error) {
	log.Printf("Refresh access token...\n")

	body := url.Values{}
	body.Set("grant_type", GrantTypeAuthCode)
	body.Set("client_id", MonzoClientID)
	body.Set("client_secret", MonzoClientSecret)
	body.Set("redirect_uri", SuccessCallbackURL)
	body.Set("code", code)

	log.Printf("\tMaking refresh token request...\n")
	resp, _ := c.h.PostForm("https://api.monzo.com/oauth2/token", body)

	if resp.StatusCode != http.StatusOK {
		baseErr := getMonzoErrorMessage(resp.StatusCode).Error()
		var e Error
		_ = json.NewDecoder(resp.Body).Decode(&e)

		log.Printf("\tFailed to make request: %s\n", e.Message)

		return nil, fmt.Errorf("%s: %s", baseErr, e.Message)
	}

	log.Printf("Access token refreshed!")

	var ac AccessToken
	_ = json.NewDecoder(resp.Body).Decode(&ac)

	return &ac, nil
}

func (c *client) RefreshAccessToken(refreshToken string) (*AccessToken, error) {
	body := url.Values{}
	body.Set("grant_type", GrantTypeRefreshToken)
	body.Set("client_id", MonzoClientID)
	body.Set("client_secret", MonzoClientSecret)
	body.Set("refresh_token", refreshToken)

	resp, _ := c.h.PostForm("https://api.monzo.com/oauth2/token", body)

	if resp.StatusCode != http.StatusOK {
		baseErr := getMonzoErrorMessage(resp.StatusCode).Error()
		var e Error
		_ = json.NewDecoder(resp.Body).Decode(&e)

		return nil, fmt.Errorf("%s: %s", baseErr, e.Message)
	}

	var ac AccessToken
	_ = json.NewDecoder(resp.Body).Decode(&ac)

	return &ac, nil
}

func (c *client) Logout(u *model.User) error {
	return nil
}

func getMonzoErrorMessage(code int) error {
	switch code {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusMethodNotAllowed:
		return ErrMethodNotAllowed
	case http.StatusNotFound:
		return ErrPageNotFound
	case http.StatusNotAcceptable:
		return ErrNotAcceptable
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		return ErrInternalServerError
	case http.StatusGatewayTimeout:
		return ErrGatewayTimeout
	default:
		return fmt.Errorf("an error occured (unrecognized status code: %d)", code)
	}
}
