package monzo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
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

const (
	// GrantTypeAuthCode is the grant required for exchanging authentication codes for access tokens.
	GrantTypeAuthCode = "authorization_code"
	// GrantTypeRefreshToken is the grant required for refreshing an access token.
	GrantTypeRefreshToken = "refresh_token"

	AccountTypeUKRetail = "uk_retail"

	// ResponseTypeCode is the response type required for the authenticating a used.
	ResponseTypeCode = "code"

	// LoginBaseURL is the base url for the Monzo login page.
	LoginBaseURL = `https://auth.monzo.com/`

	// APIBaseURL is the base url for Monzo's API.
	APIBaseURL = `https://api.monzo.com/`

	// The default HTTP request timeout.
	requestTimeoutSeconds = 5
)

// Environment Variable names.
const (
	VarClientID         = "MONZO_CLIENT_ID"
	VarClientSecret     = "MONZO_CLIENT_SECRET"
	VarOAuthCallbackURL = "OAUTH_CALLBACK_URL"
	VarWebhookURL       = "MONZO_WEBHOOK_URL"
)

// Client is a high-level interface, used to communicate and interact with
// Monzo's authentication flow and API.
type Client interface {
	// Requests an access token from Monzo, by exchanging an authorization code
	// from an earlier stage of the authentication flow.
	RequestAccessToken(code string) (*AccessToken, error)

	// Refreshes an pre-aquired access token, by using its refresh token. This
	// returns a fresh access token. Once a token has been refreshed, both the
	// original access and refresh tokens can no longer be used.
	RefreshAccessToken(refreshToken string) (*AccessToken, error)

	// Given a state token, the given HTTP response and request will be used
	// to return a redirect response to Monzo's authentication page - this
	// will be a permanent redirect to avoid backtracking (starting the flow again).
	Login(w http.ResponseWriter, r *http.Request, state string)

	// WhoAmI is used to validate a user's access token by making a request
	// to /ping/whoami. If this failed, the access token is most likely
	// invalid. Return is a *AuthenticationData instance, containing the
	// authentication status, client id and user's id.
	WhoAmI(accessToken string) (*AuthenticationData, error)

	// RegisterHook is used to register a webhook to Monz++, on the user's account.
	RegisterHook(accountID, accessToken string) error

	// Accounts returns an array of a user's personal accounts. Only a user's
	// own accounts will be returned; joint accounts will not be found.
	Accounts(accessToken string) ([]*AccountData, error)
}

var (
	defaultClient Client
)

func RequestAccessToken(code string) (*AccessToken, error) {
	return mustUseDefaultClient().RequestAccessToken(code)
}

func RefreshAccessToken(refreshToken string) (*AccessToken, error) {
	return mustUseDefaultClient().RefreshAccessToken(refreshToken)
}

func Login(w http.ResponseWriter, r *http.Request, state string) {
	mustUseDefaultClient().Login(w, r, state)
}

func WhoAmI(accessToken string) (*AuthenticationData, error) {
	return mustUseDefaultClient().WhoAmI(accessToken)
}

func RegisterHook(accountID, accessToken string) error {
	return mustUseDefaultClient().RegisterHook(accountID, accessToken)
}

func Accounts(accessToken string) ([]*AccountData, error) {
	return mustUseDefaultClient().Accounts(accessToken)
}

func mustUseDefaultClient() Client {
	if defaultClient == nil {
		defaultClient = NewClient()
	}

	return defaultClient
}

// client is an implementation of the Client interface.
type client struct {
	http *http.Client
}

// NewClient creates and returns a new instance of the Client interface.
func NewClient() Client {
	return &client{
		http: &http.Client{
			Timeout: time.Second * time.Duration(requestTimeoutSeconds),
		},
	}
}

func (c *client) RequestAccessToken(code string) (*AccessToken, error) {
	data := url.Values{
		"grant_type":    {GrantTypeAuthCode},
		"client_id":     {getEnvVar(VarClientID)},
		"client_secret": {getEnvVar(VarClientSecret)},
		"redirect_uri":  {getEnvVar(VarOAuthCallbackURL)},
		"code":          {code},
	}

	return c.makeTokenRequest(data)
}

func (c *client) RefreshAccessToken(refreshToken string) (*AccessToken, error) {
	data := url.Values{
		"grant_type":    {GrantTypeRefreshToken},
		"client_id":     {getEnvVar(VarClientID)},
		"client_secret": {getEnvVar(VarClientSecret)},
		"refresh_token": {refreshToken},
	}

	return c.makeTokenRequest(data)
}

// makeTokenRequest contains common logic which can be used to request an access token.
func (c *client) makeTokenRequest(body url.Values) (*AccessToken, error) {
	target, _ := url.Parse(APIBaseURL + "oauth2/token")
	resp, err := c.http.PostForm(target.String(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, readResponseError(resp)
	}

	var data AccessToken
	_ = json.NewDecoder(resp.Body).Decode(&data)

	return &data, nil
}

func (c *client) Login(w http.ResponseWriter, r *http.Request, state string) {
	target := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=%s&state=%s",
		LoginBaseURL,
		getEnvVar(VarClientID),
		getEnvVar(VarOAuthCallbackURL),
		ResponseTypeCode,
		state,
	)

	http.Redirect(w, r, target, http.StatusPermanentRedirect)
}

func (c *client) WhoAmI(accessToken string) (*AuthenticationData, error) {
	req, _ := http.NewRequest(http.MethodGet, path.Join(APIBaseURL, "ping/whoami"), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.http.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, readResponseError(resp)
	}

	var ad AuthenticationData
	_ = json.NewDecoder(resp.Body).Decode(&ad)

	return &ad, nil
}

// RegisterHook is used to regsiter a web hook to Monzo++, on the given account.
func (c *client) RegisterHook(accountID, accessToken string) error {
	target, _ := url.Parse(APIBaseURL + "webhooks")
	body := url.Values{
		"account_id": {accountID},
		"url":        {getEnvVar(VarWebhookURL)},
	}

	req, _ := http.NewRequest(http.MethodPost, target.String(), strings.NewReader(body.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return readResponseError(resp)
	}

	return nil
}

func (c *client) Accounts(accessToken string) ([]*AccountData, error) {
	target, _ := url.Parse(APIBaseURL + fmt.Sprintf("accounts?account_type=%s", AccountTypeUKRetail))
	req, _ := http.NewRequest(http.MethodGet, target.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, readResponseError(resp)
	}

	var list AccountList
	_ = json.NewDecoder(resp.Body).Decode(&list)

	return list.Accounts, nil
}

// reads the standard Monzo error response and returns a detailed error.
func readResponseError(resp *http.Response) error {
	var body Error
	_ = json.NewDecoder(resp.Body).Decode(&body)

	monzoMessage := getMonzoErrorMessage(resp.StatusCode)

	return fmt.Errorf("%v: %s", monzoMessage, body.Message)
}

// returns a predefined error for each status code.
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

func getEnvVar(name string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	panic(fmt.Errorf("the environment variable '%s' has not be set", name))
}
