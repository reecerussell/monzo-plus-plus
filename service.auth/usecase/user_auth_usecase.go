package usecase

import (
	"fmt"
	"os"

	"github.com/reecerussell/monzo-plus-plus/service.auth/jwt"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

// Environment variables.
var (
	PrivateKeyFilepath = os.Getenv("JWT_PRIVATE_KEY")
	ConfigFilepath     = os.Getenv("JWT_CONFIG")
)

// UserAuthUsecase is a high-level interface used to generate and
// validate JSON-Web tokens.
type UserAuthUsecase interface {
	GenerateToken(u *model.User) (string, errors.Error)
	ValidateToken(accessToken string) errors.Error
}

// userAuthUsecase is an implementation of the UserAuthUsecase interface.
type userAuthUsecase struct {
	keys   *jwt.KeyRegister
	config *jwt.Config
}

// NewUserAuthUsecase returns a new instance of the UserAuthusecase interface,
// initialised with an RSA key register, and jwt.Config from the jwt config file.
//
// This method expects the environment variables "JWT_PRIVATE_KEY" and "JWT_CONFIG",
// to be set and valid. In the event that either of which are invalid, an error
// will be returned with the relavent message.
func NewUserAuthUsecase() (UserAuthUsecase, error) {
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
	}, nil
}

// Generate token is used to generate the given user, u, an OAuth2
// access token. The token will be signed with relavent claims
// to the user, such as roles, unique identifier, etc.
//
// The user is assumed to be non-nil, therefore nil pointer errors
// will be thrown otherwise.
//
// An error will only be returned if in the process of signing
// the token fails.
func (uau *userAuthUsecase) GenerateToken(u *model.User) (string, errors.Error) {
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
		return "", errors.InternalError(fmt.Errorf("sign: %v", err))
	}

	return t.String(), nil
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
