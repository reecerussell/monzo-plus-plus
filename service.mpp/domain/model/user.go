package model

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

var (
	// the hashing algorithm used to hash the state token.
	stateTokenHash = sha256.New()

	// the encoding used to encode the state token.
	encoding = base64.URLEncoding
)

type User struct {
	id         string
	monzoID    string
	stateToken string

	token *UserToken
}

func NewUser() *User {
	u := new(User)

	id, _ := uuid.NewRandom()
	u.id = id.String()
	u.generateStateToken()

	return u
}

func NewUserFromStateToken(id, token string) *User {
	return &User{
		id:         id,
		stateToken: token,
		token:      &UserToken{},
	}
}

func NewUserFromScan() (*User, []interface{}) {
	u := new(User)
	u.token = new(UserToken)

	return u, []interface{}{
		&u.id,
		&u.monzoID,
		&u.stateToken,
		&u.token.accessToken,
		&u.token.refreshToken,
		&u.token.expires,
		&u.token.tokenType,
	}
}

func (u *User) GetID() string {
	return u.id
}

func (u *User) GetMonzoID() string {
	return u.monzoID
}

func (u *User) GetStateToken() string {
	return u.stateToken
}

func (u *User) GetToken() *UserToken {
	return u.token
}

func (u *User) SetToken(ut *UserToken) {
	u.token = ut
}

// UpdateMonzoID sets the value of the user's MonzoID.
func (u *User) UpdateMonzoID(id string) error {
	if id == "" {
		return fmt.Errorf("user: monzo id cannot be empty")
	}

	u.monzoID = id

	return nil
}

// generates a state token by hashing the user's id.
func (u *User) generateStateToken() {
	if u.id == "" {
		return
	}

	stateTokenHash.Reset()
	hashedUserID := stateTokenHash.Sum([]byte(u.id))
	u.stateToken = encoding.EncodeToString(hashedUserID)
}
