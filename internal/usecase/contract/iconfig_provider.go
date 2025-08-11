package usecasecontract

import "time"

type IConfigProvider interface {
	GetSendActivationEmail() bool
	GetAppBaseURL() string
	GetRefreshTokenExpiry() time.Duration
	GetPasswordResetTokenExpiry() time.Duration
	GetEmailVerificationTokenExpiry() time.Duration
	GetAIServiceAPIKey() string
}
