package password

import "fmt"

var (
	errEmpty                   = "password cannot be empty"
	errNotLongEnough           = "password must be at least %d characters long"
	errRequiresUppercase       = "password requires at least one uppercase character"
	errRequiresLowercase       = "password requires at least one lowercase character"
	errRequiresNonAlphanumeric = "password requires at least one alphanumeric character"
	errRequiresDigit           = "password requires at least one digit"
	errNotEnoughUniqueChars    = "password requires at least %d unique characters"
)

// Options contains a set of critrea that a
// password has to hit, to be considered valid.
type Options struct {
	// RequiredLength is the minimum length a password has to be.
	RequiredLength int

	// RequireUppercase is a flag which demands at least one
	// character to be uppercase.
	RequireUppercase bool

	// RequireLowercase is a flag which demands at least one
	// character to be lowercase.
	RequireLowercase bool

	// RequireNonAlphanumeric is a flag which demands at least one
	// character to be non alphanumeric.
	RequireNonAlphanumeric bool

	// RequireDigit is a flag which demands at least one
	// character to be a digit.
	RequireDigit bool

	// RequiredUniqueChars is an integer value, that determines
	// how many unique characters are required in a password.
	RequiredUniqueChars int
}

// Validate validates a given password against o's values. Returns
// an error is password is not valid, otherwise nil.
func (o *Options) Validate(pwd string) error {
	l := len(pwd)
	if pwd == "" || l < 1 {
		return fmt.Errorf(errEmpty)
	}

	if l < o.RequiredLength {
		return fmt.Errorf(errNotLongEnough, o.RequiredLength)
	}

	var (
		hasNonAlphanumeric bool
		hasDigit           bool
		hasLower           bool
		hasUpper           bool
		uniqueChars        []byte
	)

	for i := 0; i < l; i++ {
		// Byte value of the character.
		c := byte(pwd[i])

		if !hasNonAlphanumeric && !isLetterOrDigit(c) {
			// Character is non-alphanumeric.
			hasNonAlphanumeric = true
		}

		if !hasDigit && isDigit(c) {
			// Character is a digit.
			hasDigit = true
		}

		if !hasLower && isLower(c) {
			// Character is lowercase.
			hasLower = true
		}

		if !hasUpper && isUpper(c) {
			// Character is uppercase.
			hasUpper = true
		}

		d := true

		for _, dc := range uniqueChars {
			if dc == c {
				// Character is non-unique.
				d = false
			}
		}

		if d {
			uniqueChars = append(uniqueChars, c)
		}
	}

	if o.RequireNonAlphanumeric && !hasNonAlphanumeric {
		return fmt.Errorf(errRequiresNonAlphanumeric)
	}

	if o.RequireDigit && !hasDigit {
		return fmt.Errorf(errRequiresDigit)
	}

	if o.RequireLowercase && !hasLower {
		return fmt.Errorf(errRequiresLowercase)
	}

	if o.RequireUppercase && !hasUpper {
		return fmt.Errorf(errRequiresUppercase)
	}

	if o.RequiredUniqueChars >= 1 && len(uniqueChars) < o.RequiredUniqueChars {
		return fmt.Errorf(errNotEnoughUniqueChars, o.RequiredUniqueChars)
	}

	// Password is valid.
	return nil
}

// isDigit returns a flag indicating whether the supplied character
// is a digit - true if the character is a digit, otherwise false.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLower returns a flag indicating whether the supplied character is
// a lower case ASCII letter - true if the character is a lower case
// ASCII letter, otherwise false.
func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// isUpper returns a flag indicating whether the supplied character is
// an upper case ASCII letter - true if the character is an upper case ASCII
// letter, otherwise false.
func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// isLetterOrDigit returns a flag indicating whether the supplied character is
// an ASCII letter or digit - true if the character is an ASCII letter or digit,
// otherwise false.
func isLetterOrDigit(c byte) bool {
	return isUpper(c) || isLower(c) || isDigit(c)
}
