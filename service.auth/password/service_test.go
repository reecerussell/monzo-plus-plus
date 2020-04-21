package password

import (
	"fmt"
	"testing"
)

var (
	options = &Options{
		RequiredUniqueChars:    6,
		RequiredLength:         6,
		RequireUppercase:       true,
		RequireNonAlphanumeric: true,
		RequireLowercase:       true,
		RequireDigit:           true,
	}
)

func TestValidatePasswordWithInvalidPasswords(t *testing.T) {
	serv := NewService(options, NewHasher())

	pwd := ""
	err := serv.Validate(pwd)
	if err.Error() != errEmpty {
		t.Errorf("expected '%s' but got '%v'", errEmpty, err)
	}

	pwd = "a"
	err = serv.Validate(pwd)
	if err.Error() != fmt.Sprintf(errNotLongEnough, options.RequiredLength) {
		t.Errorf("expected '%s' but got '%v'", fmt.Sprintf(errNotLongEnough, options.RequiredLength), err)
	}

	pwd = "aaaaaa"
	err = serv.Validate(pwd)
	if err.Error() != errRequiresNonAlphanumeric {
		t.Errorf("expected '%s' but got '%v'", errRequiresNonAlphanumeric, err)
	}

	pwd = "aaaaaaa!"
	err = serv.Validate(pwd)
	if err.Error() != errRequiresDigit {
		t.Errorf("expected '%s' but got '%v'", errRequiresDigit, err)
	}

	pwd = "111111!"
	err = serv.Validate(pwd)
	if err.Error() != errRequiresLowercase {
		t.Errorf("expected '%s' but got '%v'", errRequiresLowercase, err)
	}

	pwd = "a111111!"
	err = serv.Validate(pwd)
	if err.Error() != errRequiresUppercase {
		t.Errorf("expected '%s' but got '%v'", errRequiresUppercase, err)
	}

	pwd = "aA11111!"
	err = serv.Validate(pwd)
	if err.Error() != fmt.Sprintf(errNotEnoughUniqueChars, options.RequiredUniqueChars) {
		t.Errorf("expected '%s' but got '%v'", fmt.Sprintf(errNotEnoughUniqueChars, options.RequiredUniqueChars), err)
	}
}

func TestValidatePassword(t *testing.T) {
	serv := NewService(options, NewHasher())
	pwd := "aBc_124nd!?"
	err := serv.Validate(pwd)
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
	}
}
