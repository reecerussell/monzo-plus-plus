package password

import (
	"fmt"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	pwd := "Password_1234!"

	serv := NewService(options, NewHasher())
	pwdHash, err := serv.Hash(pwd)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(pwdHash)

	ok := serv.Verify(pwd, pwdHash)
	if !ok {
		t.Errorf("verification failed")
	}
}

func TestVerifyWithInvalidHash(t *testing.T) {
	pwd := "Password_1234!"

	serv := NewService(options, NewHasher())
	pwdHash, err := serv.Hash(pwd)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(pwdHash)

	ok := serv.Verify(pwd,
		// Distort hash
		pwdHash[5:])
	if ok {
		t.Errorf("expected false but got %v", ok)
	}
}

func TestPasswordHashWithEmptyPassword(t *testing.T) {
	serv := NewService(options, NewHasher())
	pwdHash, err := serv.Hash("")
	if err == nil {
		t.Fail()
	}

	fmt.Println(pwdHash)

	ok := serv.Verify("", pwdHash)
	if ok {
		t.Errorf("should failed")
	}
}
