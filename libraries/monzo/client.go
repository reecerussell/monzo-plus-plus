package monzo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// Monzo status code errors
var (
	ErrBadRequest          = errors.New(fmt.Errorf("your request has missing arguments or is malformed"))
	ErrUnauthorized        = errors.New(fmt.Errorf("your request is not authenticated"))
	ErrForbidden           = errors.New(fmt.Errorf("your request is authenticated but has insufficent permissions"))
	ErrMethodNotAllowed    = errors.New(fmt.Errorf("you are using an incorrect HTTP verb. double check whether it should be POST/GET/DELETE/etc"))
	ErrPageNotFound        = errors.New(fmt.Errorf("the endpoint requested does not exist"))
	ErrNotAcceptable       = errors.New(fmt.Errorf("your application does not accept the content format returned according to the accpet headers sent in the request"))
	ErrTooManyRequests     = errors.New(fmt.Errorf("your application is exceeding its rate limit. back off, buddy :p"))
	ErrInternalServerError = errors.New(fmt.Errorf("something is wrong on our end. whoopsie"))
	ErrGatewayTimeout      = errors.New(fmt.Errorf("something has timed out on our end. whoopsie"))
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
	VarClientID           = "MONZO_CLIENT_ID"
	VarClientSecret       = "MONZO_CLIENT_SECRET"
	VarOAuthCallbackURL   = "OAUTH_CALLBACK_URL"
	VarWebhookURL         = "MONZO_WEBHOOK_URL"
	VarSuccessCallbackURL = "SUCCESS_CALLBACK_URL"
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

	// RegisterHook is used to register a webhook to Monzo++, on the user's account.
	RegisterHook(userID, accountID, accessToken string) error

	// Accounts returns an array of a user's personal accounts. Only a user's
	// own accounts will be returned; joint accounts will not be found.
	Accounts(accessToken string) ([]*AccountData, error)

	CreateFeedItem(accountID, accessToken, title, imageURL string, opts ...*FeedItemOpts) errors.Error
	GetBalance(accountID, accessToken string) (*Balance, errors.Error)
	GetTransactions(accountID, accessToken string, opts ...*TransactionOpts) ([]*Transaction, errors.Error)
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

// RegisterHook uses the default client to register a web hook.
func RegisterHook(userID, accountID, accessToken string) error {
	return mustUseDefaultClient().RegisterHook(userID, accountID, accessToken)
}

func Accounts(accessToken string) ([]*AccountData, error) {
	return mustUseDefaultClient().Accounts(accessToken)
}

func CreateFeedItem(accountID, accessToken, title, imageURL string, opts ...*FeedItemOpts) errors.Error {
	return mustUseDefaultClient().CreateFeedItem(accountID, accessToken, title, imageURL, opts...)
}

func GetBalance(accountID, accessToken string) (*Balance, errors.Error) {
	return mustUseDefaultClient().GetBalance(accountID, accessToken)
}

func GetTransactions(accountID, accessToken string, opts ...*TransactionOpts) ([]*Transaction, errors.Error) {
	return mustUseDefaultClient().GetTransactions(accountID, accessToken, opts...)
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
func (c *client) RegisterHook(userID, accountID, accessToken string) error {
	target, _ := url.Parse(APIBaseURL + "webhooks")
	webhookURL := fmt.Sprintf("%s?userId=%s", getEnvVar(VarWebhookURL), userID)
	body := url.Values{
		"account_id": {accountID},
		"url":        {webhookURL},
	}

	req, _ := http.NewRequest(http.MethodPost, target.String(), strings.NewReader(body.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
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

// CreateFeedItem creates a feed item in the Monzo app for the given account.
// opts is an optional parameter where you can customise the feed item - only
// the first option item will be used.
func (c *client) CreateFeedItem(accountID, accessToken, title, imageURL string, opts ...*FeedItemOpts) errors.Error {
	log.Printf("Creating Feed Item.\n")

	body := url.Values{}
	body.Add("account_id", accountID)
	body.Add("type", "basic")
	body.Add("params[title]", title)
	body.Add("params[image_url]", imageURL)

	if len(opts) > 0 {
		opt := opts[0]

		if opt.URL != "" {
			body.Add("url", opt.URL)
		}

		if opt.TitleColor != "" {
			body.Add("params[title_color]", opt.TitleColor)
		}

		if opt.Body != "" {
			body.Add("params[body]", opt.Body)
		}

		if opt.BodyColor != "" {
			body.Add("params[body_color]", opt.BodyColor)
		}

		if opt.BackgroundColor != "" {
			body.Add("params[background_color]", opt.BackgroundColor)
		}
	}

	target, _ := url.Parse(APIBaseURL + "feed")
	req, err := http.NewRequest(http.MethodPost, target.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return errors.InternalError(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return errors.InternalError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = readResponseError(resp)
		if err != nil {
			return errors.InternalError(err)
		}
	}

	log.Printf("Created Feed Item.")

	return nil
}

// FeedItemOpts is used to customise a feed item.
type FeedItemOpts struct {
	URL             string
	TitleColor      string
	Body            string
	BodyColor       string
	BackgroundColor string
}

// GetBalance requests the current balance for the given account.
func (c *client) GetBalance(accountID, accessToken string) (*Balance, errors.Error) {
	target, _ := url.Parse(APIBaseURL + fmt.Sprintf("balance?account_id=%s", accountID))
	req, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = readResponseError(resp)
		return nil, errors.InternalError(err)
	}

	var data Balance
	_ = json.NewDecoder(resp.Body).Decode(&data)

	return &data, nil
}

// Balance holds balance data for a specific account.
type Balance struct {
	Balance      int    `json:"balance"`
	TotalBalance int    `json:"total_balance"`
	Currency     string `json:"currency"`
	SpendToday   int    `json:"spend_today"`
}

// GetTransactions is used to return an array of transactions for a specific
// account. An optional parameter, opts, is used to filter the transactions
// requested - only the first option will be used.
func (c *client) GetTransactions(accountID, accessToken string, opts ...*TransactionOpts) ([]*Transaction, errors.Error) {
	qs := fmt.Sprintf("?account_id=%s", accountID)
	if len(opts) > 0 {
		opt := opts[0]

		if opt.Since != nil {
			qs += fmt.Sprintf("&since=%s", opt.Since.Format(time.RFC3339))
		}

		if opt.Before != nil {
			qs += fmt.Sprintf("&before=%s", opt.Before.Format(time.RFC3339))
		}

		if opt.Merchant {
			qs += "&expand[]=merchant"
		}

		if opt.Limit > 0 {
			qs += fmt.Sprintf("&limit=%d", opt.Limit)
		}
	}

	target, _ := url.Parse(APIBaseURL + "transactions" + qs)
	req, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = readResponseError(resp)
		return nil, errors.InternalError(err)
	}

	var list TransactionList
	_ = json.NewDecoder(resp.Body).Decode(&list)

	return list.Transactions, nil
}

// TransactionOpts is used to filter transaction records.
type TransactionOpts struct {
	// If left nil, this option will not be used. This option is nil by default.
	Since *time.Time

	// If left nil, this option will not be used. This option is nil by default.
	Before *time.Time

	// If true, the merchant field will be expanded but is false by default.
	Merchant bool

	// Limit is used to limit the number of transactions. If the limit is 0
	// this option will be ignored. Set to 0 by default.
	Limit int
}

// TransactionList is a wrapper around an array of Transation.
type TransactionList struct {
	Transactions []*Transaction `json:"transactions"`
}

// TransactionEvent is a wrapper around Transaction.
type TransactionEvent struct {
	Type string       `json:"type"`
	Data *Transaction `json:"data"`
}

// Transaction is a struct of transaction data and is used to read from Monzo.
type Transaction struct {
	AccountBalance             int                    `json:"account_balance"`
	Amount                     int                    `json:"amount"`
	Created                    time.Time              `json:"created"`
	Currency                   string                 `json:"currency"`
	Category                   string                 `json:"category"`
	Categories                 map[string]interface{} `json:"categories"`
	Description                string                 `json:"description"`
	ID                         string                 `json:"id"`
	MetaData                   map[string]interface{} `json:"metadata"`
	Notes                      string                 `json:"notes"`
	IsLoad                     bool                   `json:"is_load"`
	Settled                    *time.Time             `json:"settled"`
	Scheme                     string                 `json:"scheme"`
	LocalAmount                int                    `json:"local_amount"`
	LocalCurrency              string                 `json:"local_currency"`
	Updated                    *time.Time             `json:"updated"`
	AccountID                  string                 `json:"account_id"`
	UserID                     string                 `json:"user_id"`
	CounterParty               *CounterParty          `json:"counterparty"`
	DedupeID                   string                 `json:"dedupe_id"`
	Originator                 bool                   `json:"originator"`
	IncludeInSpending          bool                   `json:"include_in_spending"`
	CanBeExcludedFromBreakdown bool                   `json:"can_be_excluded_from_breakdown"`
	CanBeMadeSubscription      bool                   `json:"can_be_made_subscription"`
	CanSplitTheBill            bool                   `json:"can_split_the_bill"`
	CanAddToTab                bool                   `json:"can_Add_to_tab"`
	AmountIsPending            bool                   `json:"amount_is_pending"`
	Labels                     []string               `json:"labels"`
}

// Merchant is used to read merchant data from of a transactions.
type Merchant struct {
	Address         *MerchantAddress       `json:"address"`
	Created         time.Time              `json:"created"`
	GroupID         string                 `json:"group_id"`
	ID              string                 `json:"id"`
	Logo            string                 `json:"logo"`
	Emoji           string                 `json:"emoji"`
	Name            string                 `json:"name"`
	Category        string                 `json:"category"`
	Online          *bool                  `json:"online"`
	ATM             *bool                  `json:"atm"`
	MetaData        map[string]interface{} `json:"metadata"`
	DisableFeedback *bool                  `json:"disable_Feedback"`
}

// MerchantAddress holds address data for a merchant.
type MerchantAddress struct {
	Address        string  `json:"address"`
	City           string  `json:"city"`
	Country        string  `json:"country"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Postcode       string  `json:"postcode"`
	Region         string  `json:"region"`
	ShortFormatted *string `json:"short_formatted"`
	Formatted      *string `json:"formatted"`
	ZoomLevel      *int    `json:"zoom_level"`
	Approximate    *bool   `json:"approximate"`
}

// CounterParty is used to read counterparty data.
type CounterParty struct {
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	SortCode      string `json:"sort_code"`
	UserID        string `json:"user_id"`
}

// reads the standard Monzo error response and returns a detailed error.
func readResponseError(resp *http.Response) error {
	var body Error
	_ = json.NewDecoder(resp.Body).Decode(&body)

	// This is a temporary fix for while theres an issue with
	// Monzo's API. 19/04/2020
	if body.Code == "bad_request.missing_param.type" {
		return nil
	}

	monzoMessage := getMonzoErrorMessage(resp.StatusCode)

	return fmt.Errorf("%v: %s", monzoMessage, body.Message)
}

// returns a predefined error for each status code.
func getMonzoErrorMessage(code int) errors.Error {
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
		return errors.New(fmt.Errorf("an error occured (unrecognized status code: %d)", code))
	}
}

func getEnvVar(name string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	panic(fmt.Errorf("the environment variable '%s' has not be set", name))
}
