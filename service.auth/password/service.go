package password

// Service is a wrapper around the password hasher and validator.
// This keeps all password logic in the same place, making
// it easier to use and decreasing the amount of dependencies
// methods require.
type Service interface {
	Validate(pwd string) error
	Hash(pwd string) (string, error)
	Verify(pwd, pwdHash string) bool
}

// service is an implementation for the Service interface.
type service struct {
	options *Options
	hasher  *Hasher
}

// NewService is a factory method for the Service interface,
// creating an instance of an implementation for it.
func NewService(options *Options, hasher *Hasher) Service {
	return &service{
		options: options,
		hasher:  hasher,
	}
}

// Validate executes a password validation proceedure
// against a given password.
func (s *service) Validate(pwd string) error {
	return s.options.Validate(pwd)
}

// Hash takes a given password and hashes it using SHA256
// algorithm, a 256 bit key and a 64 bit sub-key.
func (s *service) Hash(pwd string) (string, error) {
	return s.hasher.Hash(pwd)
}

// Verify takes a given password and an existing password hash
// and verifies the password using the hash.
func (s *service) Verify(pwd, pwdHash string) bool {
	return s.hasher.Verify(pwd, pwdHash)
}
