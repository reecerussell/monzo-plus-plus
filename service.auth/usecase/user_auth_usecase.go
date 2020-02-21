package usecase

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/jwt"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
)

// Environment variables.
var (
	PrivateKeyFilepath = os.Getenv("JWT_PRIVATE_KEY")
	ConfigFilepath     = os.Getenv("JWT_CONFIG")
)

// UserAuthUsecase is a high-level interface used to generate and
// validate JSON-Web tokens.
type UserAuthUsecase interface {
	GenerateToken(c *dto.UserCredential) (*jwt.AccessToken, errors.Error)
	ValidateToken(accessToken string) errors.Error
	ValidateCredentials(c *dto.UserCredential) (*model.User, errors.Error)
}

// userAuthUsecase is an implementation of the UserAuthUsecase interface.
type userAuthUsecase struct {
	keys   *jwt.KeyRegister
	config *jwt.Config
	ps     password.Service
	repo   repository.UserRepository
}

// NewUserAuthUsecase returns a new instance of the UserAuthusecase interface,
// initialised with an RSA key register, and jwt.Config from the jwt config file.
//
// This method expects the environment variables "JWT_PRIVATE_KEY" and "JWT_CONFIG",
// to be set and valid. In the event that either of which are invalid, an error
// will be returned with the relavent message.
func NewUserAuthUsecase(ps password.Service, repo repository.UserRepository) (UserAuthUsecase, error) {
	keys, err := jwt.NewKeyRegisterFromFile(PrivateKeyFilepath, nil)
	if err != nil {
		return nil, err
	}

	var config jwt.Config
	err = config.LoadFromFile(ConfigFilepath)
	if err != nil {
		return nil, err
	}

	return &userAuthUsecase{
		keys:   keys,
		config: &config,
		ps:     ps,
		repo:   repo,
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
		"user_id": u.GetID(),
		"roles":   u.GetRoles(),
	}

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
	t := jwt.FromToken([]byte(accessToken))

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
