package model

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/domain"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/handler"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func init() {
	domain.RegisterEventHandler(&event.AddUserToRole{}, &handler.AddUserToRole{})
	domain.RegisterEventHandler(&event.RemoveUserFromRole{}, &handler.RemoveUserFromRole{})
	domain.RegisterEventHandler(&event.EnablePluginForUser{}, &handler.EnablePluginForUser{})
	domain.RegisterEventHandler(&event.DisablePluginForUser{}, &handler.DisablePluginForUser{})
	domain.RegisterEventHandler(&event.RegisterWebhook{}, &handler.RegisterWebhook{})
	domain.RegisterEventHandler(&event.DeleteToken{}, &handler.DeleteToken{})
}

// User is a domain model used to manage and create user, token
// records - as well as, manage user's roles.
type User struct {
	domain.Aggregate

	id           string
	username     string
	passwordHash string
	stateToken   string
	enabled      *time.Time
	accountID    *string

	roles []*Role
	token *UserToken
}

// NewUser is used to create a new User domain model. Given the data in d,
// the username and password will be set, after been validated.
func NewUser(d *dto.CreateUser, service password.Service) (*User, errors.Error) {
	id, _ := uuid.NewRandom()
	u := new(User)

	u.id = id.String()
	u.stateToken = base64.RawURLEncoding.EncodeToString([]byte(u.id))

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

// GetUsername returns the user's username.
func (u *User) GetUsername() string {
	return u.username
}

// GetStateToken returns the user's state token.
func (u *User) GetStateToken() string {
	return u.stateToken
}

// IsEnabled returns whether the user is enabled or not.
func (u *User) IsEnabled() bool {
	return u.enabled != nil
}

// HasAccount returns whether the user has specified an account or not.
func (u *User) HasAccount() bool {
	return u.accountID != nil
}

// GetRoles returns an array of the user's assigned role ids.
func (u *User) GetRoles() []string {
	ids := make([]string, len(u.roles))

	for i, r := range u.roles {
		ids[i] = r.id
	}

	return ids
}

// GetRoleNames returns an array of the user's assigned role names.
func (u *User) GetRoleNames() []string {
	names := make([]string, len(u.roles))

	for i, r := range u.roles {
		names[i] = r.name
	}

	return names
}

// GetToken returns the user's UserToken.
func (u *User) GetToken() *UserToken {
	return u.token
}

// HasValidToken returns whether the user has a valid Monzo token.
func (u *User) HasValidToken() bool {
	if u.token == nil {
		return false
	}

	if u.token.accessToken == "" {
		return false
	}

	return true
}

// Update updates the user's mutable properties, such as username.
func (u *User) Update(d *dto.UpdateUser) errors.Error {
	err := u.UpdateUsername(d.Username)
	if err != nil {
		return err
	}

	return nil
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
	err := u.VerifyPassword(currentPassword, service)
	if err != nil {
		return err
	}

	return u.setPassword(newPassword, service)
}

// VerifyPassword is used to verify a given password against the user's password
// hash. This used to given password.Service to verify the hash.
func (u *User) VerifyPassword(pwd string, service password.Service) errors.Error {
	if !service.Verify(pwd, u.passwordHash) {
		return errors.BadRequest("password is invalid")
	}

	return nil
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

	u.passwordHash = service.Hash(pwd)

	return nil
}

// UpdateToken resets the user's token data with the AccessToken given.
func (u *User) UpdateToken(d *monzo.AccessToken) {
	u.token = NewUserToken(d)
}

// ClearToken delete the user's Monzo access token data. This is typically
// used if it is discovered the token data is invalid and requires a user's
// reauthentication.
func (u *User) ClearToken() {
	u.token = nil

	u.RaiseEvent(&event.DeleteToken{
		UserID: u.id,
	})
}

// Enable marks the user as enabled, by setting the enabled date
// to the current time, in UTC form. The user cannot already
// be enabled, if it is, and arror is returned.
func (u *User) Enable() errors.Error {
	if u.enabled != nil {
		return errors.BadRequest("user already enabled")
	}

	t := time.Now().UTC()
	u.enabled = &t

	return nil
}

// AddToRole adds the user to a the given role. If the user is
// already assigned to the role, an error is returned.
func (u *User) AddToRole(r *Role) errors.Error {
	for _, ur := range u.roles {
		if ur.GetID() == r.GetID() {
			return errors.BadRequest(fmt.Sprintf("role '%s' has already been assigned", r.GetName()))
		}
	}

	u.RaiseEvent(&event.AddUserToRole{
		UserID: u.GetID(),
		RoleID: r.GetID(),
	})

	return nil
}

// RemoveFromRole unassigns the user from the given role. An error is
// returned if the user isn't already assigned.
func (u *User) RemoveFromRole(r *Role) errors.Error {
	for _, ur := range u.roles {
		if ur.GetID() == r.GetID() {
			u.RaiseEvent(&event.RemoveUserFromRole{
				UserID: u.GetID(),
				RoleID: r.GetID(),
			})

			return nil
		}
	}

	return errors.BadRequest(fmt.Sprintf("role '%s' has not already been assigned", r.GetName()))
}

// EnablePlugin is used to enable a specific plugin for the user. A
// domain event is raise to handle enabling it.
func (u *User) EnablePlugin(pluginID string) errors.Error {
	u.RaiseEvent(&event.EnablePluginForUser{
		PluginID: pluginID,
		UserID:   u.GetID(),
	})

	return nil
}

// DisablePlugin is used to disable a specific plugin for the user. A
// domain event is raise to handle disabling it.
func (u *User) DisablePlugin(pluginID string) errors.Error {
	u.RaiseEvent(&event.DisablePluginForUser{
		PluginID: pluginID,
		UserID:   u.GetID(),
	})

	return nil
}

// UpdateAccountID sets the user's Monzo account id, then raises a domain
// event to register a webhook with the given account.
func (u *User) UpdateAccountID(accountID, accessToken string) errors.Error {
	switch true {
	case accountID == "":
		return errors.BadRequest("account id cannot be empty")
	case len(accountID) > 128:
		return errors.BadRequest("account id cannot be greater than 128 characters long")
	}

	u.accountID = &accountID

	u.RaiseEvent(&event.RegisterWebhook{
		AccessToken: accessToken,
		AccountID:   accountID,
	})

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

	if u.accountID == nil {
		dm.AccountID = sql.NullString{
			Valid: false,
		}
	} else {
		dm.AccountID = sql.NullString{
			Valid:  true,
			String: *u.accountID,
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
func UserFromDataModel(dm *datamodel.User, rdm []*datamodel.Role, tdm *datamodel.UserToken) *User {
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

	if dm.AccountID.Valid {
		u.accountID = &dm.AccountID.String
	} else {
		u.accountID = nil
	}

	if rdm != nil {
		u.roles = make([]*Role, len(rdm))

		for i, r := range rdm {
			u.roles[i] = RoleFromDataModel(r)
		}
	}

	if tdm != nil {
		u.token = UserTokenFromDataModal(tdm)
	}

	return u
}

// DTO returns a data-transfer object populated with the user's data.
func (u *User) DTO() *dto.User {
	d := &dto.User{
		ID:          u.id,
		Username:    u.username,
		DateEnabled: u.enabled,
		Enabled:     u.enabled != nil,
		MonzoLinked: u.HasValidToken(),
		AccountID:   u.accountID,
	}

	if u.roles != nil {
		d.Roles = make([]*dto.Role, len(u.roles))

		for i, r := range u.roles {
			d.Roles[i] = r.DTO()
		}
	}

	return d
}
