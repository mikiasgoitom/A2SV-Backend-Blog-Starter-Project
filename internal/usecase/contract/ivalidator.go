package usecasecontract

// Validator defines the interface for generic input validation.
type IValidator interface {
	ValidateEmail(email string) error
	ValidatePasswordStrength(password string) error
}
