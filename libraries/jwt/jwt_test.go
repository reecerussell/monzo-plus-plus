package jwt

import (
	"testing"
	"time"
)

const (
	encryptedKeyData = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAtONH7jHZ+QyvrVWUkCnyNYiXvUWy8JDM0+N6zLn97jvk87hA
heYGggbGGjNUGeJni7cfPfcU1JZOaWGFc8yHuKKPKz24rOBZqFrfWEoLkRB1Y5Ww
TByBmwGYdLFjfwoRGU6G97YHDAMGEnsJzub6Z3xA32FeDjMeosgFGPlgXrQ1QMG9
QvzJovmqhe3C0Rf+jayv8PAxfV0C49EA60Z43jA++/4kbiyTzHVspTzXdbWcAK5o
1PxJySAN9GMfZXMB+mCdDzwPf7t/vH8LA4aIPBTzD+sSyQcGs0yAcJ26YfrqxBc2
QoZAdgiUH+tN9IV1uVs921IN/FjEiehf/Q8MIQIDAQABAoIBAA8ojLqVSuLoAUDR
TyXVngqGa9DcqmYmfEO1aHEHlRQFyOXzptSRtjHnR2qiqoWQx4SZz/BtaD14axHB
rmFJ3oXGeaDyByvVkS3ej6Dic52wd2XlAWUfbm0C8Te2NdRLj6tDPWQ8yNJk3nll
/ihsisdpTjZp/mvKNOMHSAYTv9OvuhKABgn1la5ijNZuPFrE6mzfi/O5TjLEhlzl
TzsDC3BqtkAnHLlcMrNUvYFbId21onCAi2hiHIEHalzaRAP4GInqqrvJ6Km6b47e
8GLuXKdRMSY4eKYNWEowtvgZcu09Dl/h2h6RgF0vla6Ra008oCQvby4t/fr6D7Do
VETZLtECgYEA7xHRALGdXb0/rsB9cLuFeGBOm5BfmMf4Ong/H4tDtk/vAhH0pcRQ
8G4vWawQmQ+lXfYxMUOl8yaeLQkS8wycmtOtppHsHDTnO4NVEcecFe3GY6jJRBFU
iJXEd4DwU1JFEXj9W0SZdjfgxCSL8w2NDWm1qcMh9KlNWZdK+lTv8U8CgYEAwbKq
O3MvVVaFESrpkl0ccZWAuK0jbpEYSNHypMNIwIeg81F0elpi0FNMj5zZTvQp90aT
CUi8eSsuYEENYfD6XVye+2jFxIJvJiS5Y/wFUKqD3S5f8vVR/L0Wng+KlM2070BG
soMR8aVENMCugq3rSqCWUaRYjUDFL+PiVHK6b48CgYEAis36Xn/BkaKGrgzfCCwo
Y2rvWQ0rJAa+zhvw8nymVNz2NWp9dB2WrPIcleD8RhI6fmKpzyOq35FFd3p0QR/K
cW6DvVmuD/Z7ZydqpgcSTf0fGc/vA5FYVnE2f/giEQ6MQYfQ+kPLIkNxQhDCHN16
TxZYPneoaGhAG8tm4g+cvKsCgYBqL+BTJxpT1lKszrGto43sVuFyXLiH7NM7WJ2E
5eoEYlZCDe2lBdgWDRRzxrD7L6+x9+azuERayOlrqz4C63DVrekaOp3w1PDKIGfi
AE+JoXEY2EuRjhucSq0TicGXQg5m+v3G70PvDYMOyHLjASV0jATcTRSem1t+nZd+
ZEcA0QKBgD1cX9gt/zrb2xWOPJ6ZdP7+y4guYIMvRRvR8S0LHX5DOr5S18ySq2RJ
IgfCxtmx4FSYu1EtyPon07uilrQYI53O2gYh/4SIdDylb4J690VyamjnhwqHvTef
BUPu315p1yzTn8ro6tX94aj44neEmlYgqvEg9WX7+I+NAPMQ4xgZ
-----END RSA PRIVATE KEY-----`
	encryptedKeyPassData = "test_pass"
)

func TestNumericTimeMapping(t *testing.T) {
	if got := NewNumericTime(time.Time{}); got != nil {
		t.Errorf("NewNumericTime from zero value got %f, want nil", *got)
	}
	if got := (*NumericTime)(nil).Time(); !got.IsZero() {
		t.Errorf("nil NumericTime got %s, want zero value", got)
	}
	if got := (*NumericTime)(nil).String(); got != "" {
		t.Errorf("nil NumericTime String got %q", got)
	}

	n := NumericTime(1234567890.12)
	d := time.Date(2009, 2, 13, 23, 31, 30, 12e7, time.UTC)

	if got := NewNumericTime(d); got == nil {
		t.Error("NewNumericTime from non-zero value got nil")
	} else if *got != n {
		t.Errorf("NewNumericTime got %f, want %f", *got, n)
	}
	if got := n.Time(); !got.Equal(d) {
		t.Errorf("Time got %s, want %s", got, d)
	}

	iso := "2009-02-13T23:31:30.12Z"
	if got := n.String(); got != iso {
		t.Errorf("String got %q, want %q", got, iso)
	}
}

func TestNewKeyRegister(t *testing.T) {
	keys, err := NewKeyRegister([]byte(encryptedKeyData), []byte(encryptedKeyPassData))
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
	}

	if keys.PrivateKey == nil {
		t.Errorf("expected key but got nil")
	}
}

func TestNewJWT(t *testing.T) {
	keys, _ := NewKeyRegister([]byte(encryptedKeyData), []byte(encryptedKeyPassData))

	n := time.Now()
	var claims Claims
	claims.Expires = NewNumericTime(n.Add(1 * time.Hour))
	claims.Issued = NewNumericTime(n)
	claims.Issuer = "jwt_test"
	claims.NotBefore = NewNumericTime(n)
	claims.Audiences = []string{"jwt_test"}

	jwt := New(&claims)
	err := jwt.Sign(keys.PrivateKey)
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
		return
	}

	ok, err := jwt.Check(keys.PublicKey)
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
		return
	}

	jwt = FromToken(jwt.data)
	ok, err = jwt.Check(keys.PublicKey)
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
		return
	}

	if !ok {
		t.Errorf("expected true but got '%v'", ok)
	}
}

func TestAddClaims(t *testing.T) {
	var claims Claims

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("wasn't expecting to panic, but received this: %v", r)
		}
	}()

	err := claims.Add("uid", "my globally unique identifier")
	if err != nil {
		t.Errorf("expected nil but got '%v'", err)
	}
}
