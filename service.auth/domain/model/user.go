package model

import (
	"regexp"
	"strings"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// User is a domain model used to manage and create user, token
// records - as well as, manage user's roles.
type User struct {
	id           string
	username     string
	passwordHash string
	stateToken   string
	enabled      *time.Time

	// TODO: add user token
	// TODO: add user's roles
}

// NewUser is used to create a new User domain model. Given the data in d,
// the username and password will be set, after been validated.
func NewUser(d *dto.CreateUser, service password.Service) (*User, errors.Error) {
	id, _ := uuid.NewRandom()
	u := new(User)

	u.id = id.String()

	// TODO: initialise nav props

	err := u.UpdateUsername(d.Username)
	if err != nil {
		return nil, err
	}

	err = u.setPassword(d.Password, service)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetID returns the user's id.
func (u *User) GetID() string {
	return u.id
}

// IsEnabled returns whether the user is enabled or not.
func (u *User) IsEnabled() bool {
	return u.enabled != nil
}

// UpdateUsername updates the user's username with the given string, username.
// A non-nil error will be returned if the username is not valid.
//
// A username must pass the validation of the following points: a username cannot be empty.
// A username must be greater than 5, but not greater than 25 characters long, and can only
// contain letters, numbers and underscores.
func (u *User) UpdateUsername(username string) errors.Error {
	username = strings.TrimSpace(username)
	l := len(username)

	if l < 1 {
		return errors.BadRequest("username is required and cannot be empty")
	}

	if l < 5 {
		return errors.BadRequest("username cannot be shorter than 5 characters")
	}

	if l > 25 {
		return errors.BadRequest("username cannot be greater than 25 characters")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9]+([_]?[a-zA-Z0-9]+)*$", username); !m {
		return errors.BadRequest("username can only contain letters, numbers and underscores. underscores cannot be at the start or end of a username.")
	}

	u.username = username

	return nil
}

// UpdatePassword is used to change a user's password. Before the password is updated,
// the currentPassword is validated against the user's current password hash. If valid,
// the newPassword is set and validated through the setPassword() function.
//
// All password validation is done using the given password.Service, which handles
// validating passwords and password hashes.
//
// An error is returned if the new or current password is in an invalid format.
func (u *User) UpdatePassword(newPassword, currentPassword string, service password.Service) errors.Error {
	if !service.Verify(currentPassword, u.passwordHash) {
		return errors.BadRequest("password is invalid")
	}

	return u.setPassword(newPassword, service)
}

// setPassword sets the user's password after hashing it using the given password.Service.
// The new password is also validated against the password options in the service.
//
// An error is only returned if the password is invalid.
func (u *User) setPassword(pwd string, service password.Service) errors.Error {
	err := service.Validate(pwd)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	ph, err := service.Hash(pwd)
	if err != nil {
		return errors.InternalError(err)
	}

	u.passwordHash = ph

	return nil
}

// DataModel returns a new instance of the User data model,
// which should be used to write data to the database.
func (u *User) DataModel() *datamodel.User {
	dm := &datamodel.User{
		ID:           u.id,
		Username:     u.username,
		StateToken:   u.stateToken,
		PasswordHash: u.passwordHash,
	}

	if u.enabled == nil {
		dm.Enabled = mysql.NullTime{
			Valid: false,
		}
	} else {
		dm.Enabled = mysql.NullTime{
			Valid: true,
			Time:  *u.enabled,
		}
	}

	return dm
}

// UserFromDataModel is used to initalise a new User domain model
// from the persistence layer. This function should only be used
// in a repository and when reading data.
//
// The data in the data model should be of that from the database.
// No data should be modified in the process of reading from the
// database and calling this method.
func UserFromDataModel(dm *datamodel.User) *User {
	u := &User{
		id:           dm.ID,
		username:     dm.Username,
		stateToken:   dm.StateToken,
		passwordHash: dm.PasswordHash,
	}

	if dm.Enabled.Valid {
		u.enabled = &dm.Enabled.Time
	} else {
		u.enabled = nil
	}

	return u
}
