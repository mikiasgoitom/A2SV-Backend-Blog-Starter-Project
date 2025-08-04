package validator

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
) // RegisterCustomValidators registers custom validation functions with the Gin validator.
func RegisterCustomValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("containsuppercase", containsUppercase)
        v.RegisterValidation("containslowercase", containsLowercase)
        v.RegisterValidation("containsdigit", containsNumber)
        v.RegisterValidation("containssymbol", containsSpecial)
        v.RegisterValidation("min", customMinLength)
        v.RegisterValidation("max", customMaxLength)
    }
}

// containsUppercase checks if the string contains at least one uppercase letter.
func containsUppercase(fl validator.FieldLevel) bool {
    for _, char := range fl.Field().String() {
        if unicode.IsUpper(char) {
            return true
        }
    }
    return false
}

// containsLowercase checks if the string contains at least one lowercase letter.
func containsLowercase(fl validator.FieldLevel) bool {
    for _, char := range fl.Field().String() {
        if unicode.IsLower(char) {
            return true
        }
    }
    return false
}

// containsNumber checks if the string contains at least one number.
func containsNumber(fl validator.FieldLevel) bool {
    for _, char := range fl.Field().String() {
        if unicode.IsNumber(char) {
            return true
        }
    }
    return false
}

// containsSpecial checks if the string contains at least one special character.
func containsSpecial(fl validator.FieldLevel) bool {
    for _, char := range fl.Field().String() {
        if strings.ContainsRune("!@#$%^&*()_+-=[]{};:'\\|,.<>/?", char) {
            return true
        }
    }
    return false
}

// customMinLength 
func customMinLength(fl validator.FieldLevel) bool {
    param := fl.Param()
    minLen := 0
    if param != "" {
        // Try to parse the param as an int
        var err error
        minLen, err = strconv.Atoi(param)
        if err != nil {
            return false
        }
    }
    return len(fl.Field().String()) >= minLen
}

// customMaxLength 
func customMaxLength(fl validator.FieldLevel) bool {
    param := fl.Param()
    maxLen := 0
    if param != "" {
        // Try to parse the param as an int
        var err error
        maxLen, err = strconv.Atoi(param)
        if err != nil {
            return false
        }
    }
    return len(fl.Field().String()) <= maxLen
}