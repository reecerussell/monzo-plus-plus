package model

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
)

var testPasswordOptions = password.DefaultOptions
var testPasswordHasher = password.NewHasher()
var testPasswordService = password.NewService(testPasswordOptions, testPasswordHasher)

func TestCreateUserWithValidUsername(t *testing.T) {
	data := &dto.CreateUser{Password: "test-Password1"}

	usernames := []string{
		"test_username",
		"test_user1",
		"helloworld",
	}

	for _, un := range usernames {
		data.Username = un
		_, err := NewUser(data, testPasswordService)
		if err != nil {
			t.Errorf("expected nil but got %v", err.Text())
			return
		}
	}
}

func TestCreateUserWithInvalidUsername(t *testing.T) {
	data := &dto.CreateUser{
		Username: "",
		Password: "test-Password1",
	}

	// Test empty username.
	_, err := NewUser(data, testPasswordService)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	// Test short username.
	data.Username = "test"
	_, err = NewUser(data, testPasswordService)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	// Test long username.
	data.Username = "a_super_long_username_longer_than_25_chars"
	_, err = NewUser(data, testPasswordService)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	// Test invalid usernames.
	usernames := []string{
		"test username",
		"my.test.username",
		"#myusername",
	}

	for _, un := range usernames {
		data.Username = un
		_, err = NewUser(data, testPasswordService)
		if err == nil {
			t.Errorf("expected an error but got nil")
			return
		}
	}
}

func TestCreateUserWithValidPassword(t *testing.T) {
	data := &dto.CreateUser{
		Username: "test_user",
		Password: "my-testPassword1",
	}

	_, err := NewUser(data, testPasswordService)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
	}
}

func TestCreateUserWithInvalidPassword(t *testing.T) {
	data := &dto.CreateUser{
		Username: "test_user",
		Password: "test",
	}

	_, err := NewUser(data, testPasswordService)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
}

func TestUserGetMethods(t *testing.T) {
	u, err := createUser("hello_user", "my-testPass1")
	if err != nil {
		t.Errorf("expected nil but got error: %s", err.Text())
		return
	}

	if u.GetID() != u.id {
		t.Errorf("expected %s but got %s", u.id, u.GetID())
		return
	}

	if u.GetUsername() != u.username {
		t.Errorf("expected %s but got %s", u.username, u.GetUsername())
		return
	}

	if u.GetStateToken() != u.stateToken {
		t.Errorf("expected %s but got %s", u.stateToken, u.GetStateToken())
		return
	}

	if u.IsEnabled() != (u.enabled != nil) {
		t.Errorf("expected %v but got %v", (u.enabled != nil), u.IsEnabled())
		return
	}

	if u.HasAccount() != (u.accountID != nil) {
		t.Errorf("expected %v but got %v", (u.accountID != nil), u.HasAccount())
		return
	}

	// Test Get Role Methods

	r, err := NewRole(&dto.CreateRole{Name: "role"})
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	err = u.AddToRole(r)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	roles := u.GetRoles()
	if len(roles) < 1 {
		t.Errorf("expected 1 role but got none")
		return
	}

	if roles[0] != u.roles[0].id {
		t.Errorf("expected %s but got %s", u.roles[0].id, roles[0])
		return
	}

	roles = u.GetRoleNames()
	if len(roles) < 1 {
		t.Errorf("expected 1 role but got none")
		return
	}

	if roles[0] != u.roles[0].name {
		t.Errorf("expected %s but got %s", u.roles[0].name, roles[0])
		return
	}

	if u.GetToken() != u.token {
		t.Errorf("expected %v but got %v", u.token, u.GetToken())
		return
	}
}

func TestHasValidToken(t *testing.T) {
	u, err := createUser("hello_user", "my-testPass1")
	if err != nil {
		t.Errorf("expected nil but got error: %s", err.Text())
		return
	}

	// nil token
	u.token = nil
	if u.HasValidToken() {
		t.Errorf("expected false but got true")
		return
	}

	data := &monzo.AccessToken{
		AccessToken:  "", // empty token
		TokenType:    "bearer",
		ExpiresIn:    26300,
		RefreshToken: "secure oauth refresh token",
	}

	u.UpdateToken(data)

	if u.HasValidToken() {
		t.Errorf("expected false but got true")
		return
	}

	data.AccessToken = "valid token"
	u.UpdateToken(data)

	if !u.HasValidToken() {
		t.Errorf("expected true but got false")
		return
	}
}

func TestUpdateUser(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	data := &dto.UpdateUser{
		ID:       "some random UUID",
		Username: "new_test_username",
	}

	err = u.Update(data)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	// Invalid username
	data.Username = "random invalid username"
	err = u.Update(data)
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
}

func TestUpdatePassword(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	// Invalid password
	err = u.UpdatePassword("some random password", "no the current password", testPasswordService)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	// Valid password
	err = u.UpdatePassword("my-new-test-Password1", "test-Password1", testPasswordService)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}
}

func TestUpdateToken(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	data := &monzo.AccessToken{
		AccessToken:  "secure oauth token",
		TokenType:    "bearer",
		ExpiresIn:    26300,
		RefreshToken: "secure oauth refresh token",
	}

	u.UpdateToken(data)

	ut := u.GetToken()

	if ut.accessToken != data.AccessToken {
		t.Errorf("expected access token to be '%s'", data.AccessToken)
		return
	}

	if ut.tokenType != data.TokenType {
		t.Errorf("expected token type to be '%s'", data.TokenType)
		return
	}

	if ut.refreshToken != data.RefreshToken {
		t.Errorf("expected refresh token to be '%s'", data.RefreshToken)
		return
	}
}

func TestClearToken(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	// set token
	data := &monzo.AccessToken{
		AccessToken:  "secure oauth token",
		TokenType:    "bearer",
		ExpiresIn:    26300,
		RefreshToken: "secure oauth refresh token",
	}

	u.UpdateToken(data)

	c := len(u.GetRaisedEvents())
	u.ClearToken()

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected domain event to be raised")
		return
	}

	if u.token != nil {
		t.Errorf("expected user token to be nil")
		return
	}
}

func TestEnableUser(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	err = u.Enable()
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	err = u.Enable()
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestAddToRole(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	r, err := NewRole(&dto.CreateRole{Name: "test_tole"})
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	c := len(u.GetRaisedEvents())
	err = u.AddToRole(r)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected a domain event to be raised")
		return
	}

	err = u.AddToRole(r)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}
}

func TestRemoveRole(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	r, err := NewRole(&dto.CreateRole{Name: "test_tole"})
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	err = u.RemoveFromRole(r)
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	err = u.AddToRole(r)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	c := len(u.GetRaisedEvents())
	err = u.RemoveFromRole(r)
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected a domain event to be raised")
		return
	}
}

func TestEnablePlugin(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	c := len(u.GetRaisedEvents())
	u.EnablePlugin("test plugin id")

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected a domain event to be raised")
		return
	}
}

func TestDisablePlugin(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	c := len(u.GetRaisedEvents())
	u.DisablePlugin("test plugin id")

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected a domain event to be raised")
		return
	}
}

func TestUpdateAccountID(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	err = u.UpdateAccountID("", "access-token")
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	longTestAccountID := `hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh
							hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh
							hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh
							hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh`

	err = u.UpdateAccountID(longTestAccountID, "access-token")
	if err == nil {
		t.Errorf("expected an error but got nil")
		return
	}

	c := len(u.GetRaisedEvents())
	err = u.UpdateAccountID("account-id", "access-token")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	if c+1 != len(u.GetRaisedEvents()) {
		t.Errorf("expected a domain event to be raised")
		return
	}
}

func TestUserDataModel(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	dm := u.DataModel()

	if dm.ID != u.id {
		t.Errorf("expected %s but got %s", u.id, dm.ID)
		return
	}

	if dm.StateToken != u.stateToken {
		t.Errorf("expected %s but got %s", u.stateToken, dm.StateToken)
		return
	}

	if dm.Username != u.username {
		t.Errorf("expected %s but got %s", u.username, dm.Username)
		return
	}

	if dm.PasswordHash != u.passwordHash {
		t.Errorf("expected %s but got %s", u.passwordHash, dm.PasswordHash)
		return
	}

	// enabled == nil
	if dm.Enabled.Valid {
		t.Errorf("enabled was expected to be invalid")
		return
	}

	// accountID == nil
	if dm.AccountID.Valid {
		t.Errorf("account id was expected to be invalid")
		return
	}

	u.Enable()
	u.UpdateAccountID("hello", "access token")
	dm = u.DataModel()

	if !dm.Enabled.Valid {
		t.Errorf("enabled was expected to be valid")
	}

	if !dm.AccountID.Valid || dm.AccountID.String != "hello" {
		t.Errorf("account id was expected to be valid and have a value of 'hello'")
		return
	}
}

func TestUserFromDataModel(t *testing.T) {
	udm := &datamodel.User{
		ID:           "id",
		StateToken:   "state token",
		Username:     "username",
		PasswordHash: "password hash",
		Enabled: mysql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
		AccountID: sql.NullString{
			Valid:  true,
			String: "account id",
		},
	}
	tdm := &datamodel.UserToken{
		AccessToken:  "access token",
		RefreshToken: "refresh token",
		TokenType:    "token type",
		Expires:      time.Now(),
	}
	rdm := []*datamodel.Role{
		&datamodel.Role{
			ID:   "id",
			Name: "name",
		},
	}

	u := UserFromDataModel(udm, rdm, tdm)

	if udm.ID != u.id {
		t.Errorf("expected %s but got %s", u.id, udm.ID)
		return
	}

	if udm.StateToken != u.stateToken {
		t.Errorf("expected %s but got %s", u.stateToken, udm.StateToken)
		return
	}

	if udm.Username != u.username {
		t.Errorf("expected %s but got %s", u.username, udm.Username)
		return
	}

	if udm.PasswordHash != u.passwordHash {
		t.Errorf("expected %s but got %s", u.passwordHash, udm.PasswordHash)
		return
	}

	if udm.AccountID.String != *u.accountID {
		t.Errorf("expected %s but got %s", *u.accountID, udm.AccountID.String)
		return
	}

	if udm.Enabled.Time != *u.enabled {
		t.Errorf("expected %v but got %v", *u.enabled, udm.Enabled.Time)
		return
	}

	if tdm.AccessToken != u.token.accessToken {
		t.Errorf("expected %s but got %s", u.token.accessToken, tdm.AccessToken)
		return
	}

	if tdm.RefreshToken != u.token.refreshToken {
		t.Errorf("expected %s but got %s", u.token.refreshToken, tdm.RefreshToken)
		return
	}

	if tdm.TokenType != u.token.tokenType {
		t.Errorf("expected %s but got %s", u.token.tokenType, tdm.TokenType)
		return
	}

	if tdm.Expires != u.token.expires {
		t.Errorf("expected %v but got %v", u.token.expires, tdm.Expires)
		return
	}

	if len(u.roles) != 1 {
		t.Errorf("expected 1 role but got %d", len(u.roles))
		return
	}

	r := u.roles[0]

	if rdm[0].ID != r.id {
		t.Errorf("expected %s but got %s", rdm[0].ID, r.id)
		return
	}

	if rdm[0].Name != r.name {
		t.Errorf("expected %s but got %s", rdm[0].Name, r.name)
		return
	}

	udm.AccountID = sql.NullString{Valid: false}
	udm.Enabled = mysql.NullTime{Valid: false}
	u = UserFromDataModel(udm, nil, nil)

	if u.accountID != nil {
		t.Errorf("expected nil but got %s", *u.accountID)
		return
	}

	if u.enabled != nil {
		t.Errorf("expected nil bot got %v", *u.enabled)
		return
	}

	if u.token != nil {
		t.Errorf("expected nil")
		return
	}

	if u.roles != nil {
		t.Errorf("expected nil")
		return
	}
}

func TestUserDTO(t *testing.T) {
	u, err := createUser("test_user", "test-Password1")
	if err != nil {
		t.Errorf("expected nil but got %s", err.Text())
		return
	}

	d := u.DTO()

	if d.ID != u.id {
		t.Errorf("expected %s but got %s", d.ID, u.id)
		return
	}

	if d.Username != u.username {
		t.Errorf("expected %s but got %s", d.Username, u.username)
		return
	}

	if d.DateEnabled != u.enabled {
		t.Errorf("expected %v but got %v", d.DateEnabled, u.enabled)
		return
	}

	if d.Enabled != (u.enabled != nil) {
		t.Errorf("expected %v but got %v", (u.enabled != nil), d.Enabled)
		return
	}

	if d.MonzoLinked != (u.token != nil) {
		t.Errorf("expected %v but got %v", (u.token != nil), d.MonzoLinked)
		return
	}

	if d.AccountID != u.accountID {
		t.Errorf("expected %v but got %v", d.AccountID, u.accountID)
		return
	}

	if u.roles != nil {
		t.Errorf("expected nil")
		return
	}

	u.roles = []*Role{
		&Role{
			id:   "id",
			name: "name",
		},
	}

	d = u.DTO()

	if len(u.roles) != 1 {
		t.Errorf("expected 1 role but got %d", len(u.roles))
		return
	}

	r := u.roles[0]

	if d.Roles[0].ID != r.id {
		t.Errorf("expected %s but got %s", d.Roles[0].ID, r.id)
		return
	}

	if d.Roles[0].Name != r.name {
		t.Errorf("expected %s but got %s", d.Roles[0].Name, r.name)
		return
	}
}

func createUser(u, p string) (*User, errors.Error) {
	return NewUser(&dto.CreateUser{
		Username: u,
		Password: p,
	}, testPasswordService)
}
