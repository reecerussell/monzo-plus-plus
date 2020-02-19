package password

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 64
	keySize    = 256
	iterations = 15000
)

// Hasher provides methods to hash and verify passwords,
// using a 64 bit salt and a 256 bit key, iterating 15000 times.
type Hasher struct {
}

// NewHasher is a factory to create an instance of Hasher.
func NewHasher() *Hasher {
	return &Hasher{}
}

// Hash takes a password and hashes it using a SHA256 algorithm.
// The password hash is provided in a format of <hash>.<salt>
func (*Hasher) Hash(pwd string) (string, error) {
	if pwd == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	saltBytes := make([]byte, saltSize)
	rand.Read(saltBytes)

	saltString := base64.StdEncoding.EncodeToString(saltBytes)
	salt := bytes.NewBufferString(saltString).Bytes()

	df := pbkdf2.Key([]byte(pwd), salt, iterations, keySize, sha256.New)

	cipherText := base64.StdEncoding.EncodeToString(df)

	return fmt.Sprintf("%s.%s", cipherText, saltString), nil
}

// Verify validates a password with a given hash. Hashes the given
// password then compares it with the existing hash.
func (*Hasher) Verify(pwd, pwdHash string) bool {
	if pwd == "" || pwdHash == "" || !strings.Contains(pwdHash, ".") {
		return false
	}

	cipherText := strings.Split(pwdHash, ".")[0]
	saltString := strings.Split(pwdHash, ".")[1]

	saltBytes := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(pwd), saltBytes, iterations, keySize, sha256.New)
	newCipherText := base64.StdEncoding.EncodeToString(df)

	return newCipherText == cipherText
}
