package usecase

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.auth/jwt"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
)

// Environment variables.
var (
	PrivateKeyFilepath = os.Getenv("JWT_PRIVATE_KEY")
	PrivateKeyPassword = os.Getenv("JWT_PRIVATE_KEY_PASS")
	ConfigFilepath     = os.Getenv("JWT_CONFIG")
)

// UserAuthUsecase is a high-level interface used to generate and
// validate JSON-Web tokens.
type UserAuthUsecase interface {
	GenerateToken(c *dto.UserCredential) (*jwt.AccessToken, errors.Error)
	ValidateToken(accessToken string) errors.Error
	ValidateCredentials(c *dto.UserCredential) (*model.User, errors.Error)
	WithUser(ctx context.Context, accessToken string) (context.Context, errors.Error)
	Login(code, stateToken string) errors.Error
	Register(d *dto.CreateUser) (string, errors.Error)
	GetStateToken(id string) (string, errors.Error)
	GetMonzoAccessToken(id string) (string, errors.Error)
}

// userAuthUsecase is an implementation of the UserAuthUsecase interface.
type userAuthUsecase struct {
	keys   *jwt.KeyRegister
	config *jwt.Config
	ps     password.Service
	repo   repository.UserRepository
	serv   *service.UserService
}

// NewUserAuthUsecase returns a new instance of the UserAuthusecase interface,
// initialised with an RSA key register, and jwt.Config from the jwt config file.
//
// This method expects the environment variables "JWT_PRIVATE_KEY" and "JWT_CONFIG",
// to be set and valid. In the event that either of which are invalid, an error
// will be returned with the relavent message.
func NewUserAuthUsecase(ps password.Service, repo repository.UserRepository, serv *service.UserService) (UserAuthUsecase, error) {
	keys, err := jwt.NewKeyRegisterFromFile(PrivateKeyFilepath, []byte(PrivateKeyPassword))
	if err != nil {
		return nil, fmt.Errorf("jwt: key: %v", err)
	}

	var config jwt.Config
	err = config.LoadFromFile(ConfigFilepath)
	if err != nil {
		return nil, fmt.Errorf("jwt: config: %v", err)
	}

	return &userAuthUsecase{
		keys:   keys,
		config: &config,
		ps:     ps,
		repo:   repo,
		serv:   serv,
	}, nil
}

// GenerateToken is used to generate a JSON-Web token for a user
// with the given credentials. The credentials are firstly verified,
// then if valid, a token will be generated.
//
// If the credentials are not valid, a fairly non-descriptive error
// will be returned, so that an end-user doesn't know whether the username
// or password is invalid - or both!
//
// An error could be returned if, at any point in the process, the credentials
// could not be verified, or if the token could not be generated.
func (uau *userAuthUsecase) GenerateToken(c *dto.UserCredential) (*jwt.AccessToken, errors.Error) {
	u, err := uau.ValidateCredentials(c)
	if err != nil {
		return nil, err
	}

	return uau.generateToken(u)
}

// generateToken is used to generate the given user, u, an OAuth2
// access token. The token will be signed with relavent claims
// to the user, such as roles, unique identifier, etc.
//
// The user is assumed to be non-nil, therefore nil pointer errors
// will be thrown otherwise.
//
// An error will only be returned if in the process of signing
// the token fails.
func (uau *userAuthUsecase) generateToken(u *model.User) (*jwt.AccessToken, errors.Error) {
	c := new(jwt.Claims)
	c.Set = map[string]interface{}{
		"user_id":     u.GetID(),
		"username":    u.GetUsername(),
		"enabled":     u.IsEnabled(),
		"roles":       u.GetRoleNames(),
		"has_account": u.HasAccount(),
	}

	exp := time.Now().UTC().Add(time.Duration(uau.config.ExpiryHours) * time.Hour)
	c.Expires = jwt.NewNumericTime(exp)

	t := jwt.New(c)
	t.Config = uau.config

	// An error will not be returned as the config has been pre-set, above.
	_ = t.LoadClaimsFromConfig()

	err := t.Sign(uau.keys.PrivateKey)
	if err != nil {
		return nil, errors.InternalError(fmt.Errorf("sign: %v", err))
	}

	return t.AccessToken(), nil
}

// ValidateToken is used to validate the signature of a given JSON-Web token.
//
// An error is returned if the token is invalid, or cannot be checked.
func (uau *userAuthUsecase) ValidateToken(accessToken string) errors.Error {
	t, err := jwt.FromToken([]byte(accessToken))
	if err != nil {
		return errors.Unauthorised(fmt.Sprintf("invalid token: %v", err.Error()))
	}

	ok, err := t.Check(uau.keys.PublicKey)
	if err != nil {
		return errors.Unauthorised(fmt.Sprintf("invalid token: %v", err.Error()))
	}

	if !ok {
		return errors.Unauthorised("invalid token")
	}

	return nil
}

// ValidateCredentials is used to validate a username and password, given as c.
// A non-descriptive is returned if the credentials could not be validated. If the
// credentials are valid, the user in which they belong to, is returned.
func (uau *userAuthUsecase) ValidateCredentials(c *dto.UserCredential) (*model.User, errors.Error) {
	u, err := uau.repo.GetByUsername(c.Username)
	if err != nil {
		if err.ErrorCode() == http.StatusNotFound {
			return nil, errors.BadRequest("the username and/or password is invalid")
		}

		return nil, err
	}

	err = u.VerifyPassword(c.Password, uau.ps)
	if err != nil {
		return nil, errors.BadRequest("the username and/or password is invalid")
	}

	return u, nil
}

// WithUser is used to propogate the given context with a user's
// id and domain object using the given access token. If the accessToken is
// invalid, an error will be returned with an appropriate message.
func (uau *userAuthUsecase) WithUser(ctx context.Context, accessToken string) (context.Context, errors.Error) {
	token, tErr := jwt.FromToken([]byte(accessToken))
	if tErr != nil {
		return nil, errors.Unauthorised(tErr.Error())
	}

	valid, validationError := token.Check(uau.keys.PublicKey)
	if !valid {
		if validationError != nil {
			return ctx, errors.Unauthorised(validationError.Error())
		}

		return ctx, errors.Unauthorised("invalid token")
	}

	userID, ok := token.Claims.String(jwt.ClaimUserID)
	if !ok {
		return ctx, errors.Unauthorised("user id not found in jwt payload")
	}

	u, err := uau.repo.Get(userID)
	if err != nil {
		return ctx, errors.Unauthorised("user not found")
	}

	ctx = context.WithValue(ctx, util.ContextKey("user"), u)
	ctx = context.WithValue(ctx, util.ContextKey("user_id"), u.GetID())

	return ctx, nil
}

// Login is used to request an access token from Monzo, for the user
// with the given state token. Used in the authenticate flow to
// complete the Monzo authentication process.
func (uau *userAuthUsecase) Login(code, stateToken string) errors.Error {
	u, err := uau.repo.GetByStateToken(stateToken)
	if err != nil {
		return nil
	}

	ac, tErr := monzo.RequestAccessToken(code)
	if tErr != nil {
		return errors.InternalError(tErr)
	}

	u.UpdateToken(ac)

	err = uau.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

// Register is used create a new user then return the newly created user's
// state token. This is used by the frontend, to register a user then
// used the state token to state the Monzo authentication flow.
func (uau *userAuthUsecase) Register(d *dto.CreateUser) (string, errors.Error) {
	u, err := model.NewUser(d, uau.ps)
	if err != nil {
		return "", err
	}

	err = uau.serv.ValidateUsername(u)
	if err != nil {
		return "", err
	}

	err = uau.repo.Insert(u)
	if err != nil {
		return "", err
	}

	return u.GetStateToken(), nil
}

// GetStateToken returns a user's state token. An error is returned
// if the user does not exist.
func (uau *userAuthUsecase) GetStateToken(id string) (string, errors.Error) {
	u, err := uau.repo.Get(id)
	if err != nil {
		return "", err
	}

	return u.GetStateToken(), nil
}

// GetMonzoAccessToken is used to get a user's Monzo access token. If the user does
// not have a Monzo access token an error is returned. However, if they do and it's
// expired, the refresh token will be used to request a new one, which is then saved.
func (uau *userAuthUsecase) GetMonzoAccessToken(id string) (string, errors.Error) {
	u, err := uau.repo.Get(id)
	if err != nil {
		return "", err
	}

	td := u.GetToken()
	if td == nil {
		return "", errors.BadRequest("user is not linked to monzo")
	}

	accessToken := td.GetAccessToken()

	if time.Now().UTC().Unix() > td.GetExpiryDate().UTC().Unix() {
		ac, err := monzo.RefreshAccessToken(td.GetRefreshToken())
		if err != nil {
			u.ClearToken()

			if err := uau.repo.Update(u); err != nil {
				return "", err
			}

			return "", errors.InternalError(fmt.Errorf("failed to refresh access token: %v", err))
		}

		u.UpdateToken(ac)

		if err := uau.repo.Update(u); err != nil {
			return "", err
		}

		accessToken = ac.AccessToken
	}

	return accessToken, nil
}
