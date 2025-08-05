package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// UserUsecase implements the UserUseCase interface.
type UserUsecase struct {
	userRepo                   UserRepository
	tokenRepo                  TokenRepository
	emailVerificationTokenRepo EmailVerificationTokenRepository
	hasher                     Hasher
	jwtService                 JWTService
	mailService                MailService
	logger                     AppLogger
	cfg                        ConfigProvider
	validator                  Validator
	uuidGenerator              UUIDGenerator
}

// NewUserUsecase creates a new UserUsecase instance.
func NewUserUsecase(
	userRepo UserRepository,
	tokenRepo TokenRepository,
	emailVerificationTokenRepo EmailVerificationTokenRepository,
	hasher Hasher,
	jwtService JWTService,
	mailService MailService,
	logger AppLogger,
	cfg ConfigProvider,
	validator Validator,
	uuidGenerator UUIDGenerator,
) *UserUsecase {
	return &UserUsecase{
		userRepo:                   userRepo,
		tokenRepo:                  tokenRepo,
		emailVerificationTokenRepo: emailVerificationTokenRepo,
		hasher:                     hasher,
		jwtService:                 jwtService,
		mailService:                mailService,
		logger:                     logger,
		cfg:                        cfg,
		validator:                  validator,
		uuidGenerator:              uuidGenerator,
	}
}

// check if UserUseCase implements the IUserUseCase
var _ IUserUseCase = (*UserUsecase)(nil)

// Register handles user registration.
func (uc *UserUsecase) Register(ctx context.Context, username, email, password, firstName, lastName string) (*entity.User, error) {
	// Validate input fields using the injected validator
	if err := uc.validator.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email format: %w", err)
	}
	if err := uc.validator.ValidatePasswordStrength(password); err != nil {
		return nil, fmt.Errorf("weak password: %w", err)
	}

	// Check if user with same username or email already exists
	const errUserNotFound = "not found"
	existingUserByEmail, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil && err.Error() != errUserNotFound {
		uc.logger.Errorf("failed to check for existing user by email: %v", err)
		return nil, fmt.Errorf("internal server error")
	}
	if existingUserByEmail != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	existingUserByUsername, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil && err.Error() != errUserNotFound {
		uc.logger.Errorf("failed to check for existing user by username: %v", err)
		return nil, fmt.Errorf("internal server error")
	}
	if existingUserByUsername != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Hash the password
	hashedPassword, err := uc.hasher.HashPassword(password)
	if err != nil {
		uc.logger.Errorf("failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	// Initialize firstName and lastName as pointers, setting to nil if empty
	var pFirstName *string
	if firstName != "" {
		pFirstName = &firstName
	}
	var pLastName *string
	if lastName != "" {
		pLastName = &lastName
	}

	// Create new user entity, initializing new fields to their zero values or nil
	user := &entity.User{
		ID:            uc.uuidGenerator.NewUUID(),
		Username:      username,
		Email:         email,
		PasswordHash:  hashedPassword,
		Role:          entity.UserRoleUser,
		PackageID:     nil,
		PackageExpiry: nil,
		IsActive:      false,
		AvatarURL:     nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		FirstName:     pFirstName,
		LastName:      pLastName,
	}

	// Save user to database
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		uc.logger.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to register user")
	}

	// Send activation email if required, using config from injected ConfigProvider
	if uc.cfg.GetSendActivationEmail() {
		// Generate email verification token
		emailVerificationTokenString, err := uc.jwtService.GenerateEmailVerificationToken(user.ID)
		if err != nil {
			uc.logger.Errorf("failed to generate email verification token for user %s: %v", user.ID, err)
		} else {
			// Hash the token before storing it
			hashedEmailVerificationToken := uc.hasher.HashString(emailVerificationTokenString)

			emailTokenEntity := &entity.EmailVerificationToken{
				ID:        uc.uuidGenerator.NewUUID(),
				UserID:    user.ID,
				TokenHash: hashedEmailVerificationToken,
				ExpiresAt: time.Now().Add(uc.cfg.GetEmailVerificationTokenExpiry()),
				Used:      false,
				CreatedAt: time.Now(),
			}

			if err := uc.emailVerificationTokenRepo.CreateEmailVerificationToken(ctx, emailTokenEntity); err != nil {
				uc.logger.Errorf("failed to store email verification token for user %s: %v", user.ID, err)
				// Log but don't fail registration if token storage fails
			} else {
				activationLink := fmt.Sprintf("%s/verify-email?token=%s", uc.cfg.GetAppBaseURL(), emailVerificationTokenString)
				if err := uc.mailService.SendActivationEmail(user.Email, user.Username, activationLink); err != nil {
					uc.logger.Warnf("failed to send activation email to %s: %v", user.Email, err)
				}
			}
		}
	}

	return user, nil
}

// Login handles user login and token generation.
func (uc *UserUsecase) Login(ctx context.Context, login, password string) (*entity.User, string, string, error) {
	// Retrieve user by username or email
	var user *entity.User
	var err error

	if uc.validator.ValidateEmail(login) == nil {
		user, err = uc.userRepo.GetUserByEmail(ctx, login)
	} else {
		user, err = uc.userRepo.GetUserByUsername(ctx, login)
	}

	if err != nil {
		const errUserNotFound = "not found"
		if err.Error() == errUserNotFound {
			return nil, "", "", errors.New("invalid credentials")
		}
		uc.logger.Errorf("failed to retrieve user for login: %v", err)
		return nil, "", "", errors.New("internal server error")
	}

	// Check if the user's email is active/verified
	if !user.IsActive {
		return nil, "", "", errors.New("account not active. Please verify your email")
	}

	// Verify password
	if !uc.hasher.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Generate access and refresh tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		uc.logger.Errorf("failed to generate access token: %v", err)
		return nil, "", "", errors.New("failed to generate token")
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		uc.logger.Errorf("failed to generate refresh token: %v", err)
		return nil, "", "", errors.New("failed to generate token")
	}

	// Create token entity with all fields from the schema
	tokenEntity := &entity.Token{
		ID:        uc.uuidGenerator.NewUUID(),
		UserID:    user.ID,
		TokenType: entity.TokenTypeRefresh,
		TokenHash: uc.hasher.HashString(refreshToken),
		ExpiresAt: time.Now().Add(uc.cfg.GetRefreshTokenExpiry()),
		CreatedAt: time.Now(),
		Revoke:    false,
	}
	if err := uc.tokenRepo.CreateToken(ctx, tokenEntity); err != nil {
		uc.logger.Errorf("failed to store refresh token for user %s: %v", user.ID, err)
		return nil, "", "", errors.New("failed to store token")
	}

	return user, accessToken, refreshToken, nil
}

// Authenticate handles user authentication using access tokens.
func (uc *UserUsecase) Authenticate(ctx context.Context, accessToken string) (*entity.User, error) {
	claims, err := uc.jwtService.ParseAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	const errUserNotFound = "not found" // Changed to camelCase
	user, err := uc.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errUserNotFound {
			return nil, errors.New("user not found")
		}
		uc.logger.Errorf("failed to retrieve user during authentication: %v", err)
		return nil, errors.New("internal server error")
	}

	return user, nil
}

// RefreshToken handles refreshing expired access tokens using refresh tokens.
func (uc *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Parse the refresh token to get the user claims.
	claims, err := uc.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	const errTokenNotFound = "not found" // Placeholder
	// Retrieve the stored token using the UserID from the claims.
	storedToken, err := uc.tokenRepo.GetTokenByUserID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errTokenNotFound {
			return "", "", errors.New("refresh token not found or invalidated, please log in again")
		}
		uc.logger.Errorf("failed to retrieve stored refresh token: %v", err)
		return "", "", errors.New("internal server error")
	}

	// Check if the token has been revoked.
	if storedToken.Revoke {
		return "", "", errors.New("refresh token has been revoked, please log in again")
	}

	// Validate refresh token against the stored hash.
	if !uc.hasher.CheckHash(refreshToken, storedToken.TokenHash) {
		uc.logger.Warnf("refresh token mismatch for user %s", claims.UserID)
		_ = uc.tokenRepo.DeleteToken(ctx, storedToken.ID) // Invalidate the stored token
		return "", "", errors.New("invalid refresh token")
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		// Refresh token expired
		_ = uc.tokenRepo.DeleteToken(ctx, storedToken.ID) // Delete the expired token
		return "", "", errors.New("refresh token expired, please log in again")
	}

	// Generate new access token
	newAccessToken, err := uc.jwtService.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		uc.logger.Errorf("failed to generate new access token during refresh: %v", err)
		return "", "", errors.New("failed to generate new access token")
	}

	// Generate a new refresh token
	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		uc.logger.Errorf("failed to generate new refresh token during refresh: %v", err)
		return "", "", errors.New("failed to generate new refresh token")
	}

	// Hash the new refresh token before storing it in the database.
	newHashedRefreshToken := uc.hasher.HashString(newRefreshToken)

	// Update the stored refresh token with the new hash and expiry.
	err = uc.tokenRepo.UpdateToken(ctx, storedToken.ID, newHashedRefreshToken, time.Now().Add(uc.cfg.GetRefreshTokenExpiry()))
	if err != nil {
		uc.logger.Errorf("failed to update refresh token in db: %v", err)
		return "", "", errors.New("failed to update token")
	}

	// Return both the new access token and the new refresh token.
	return newAccessToken, newRefreshToken, nil
}

// ForgotPassword handles the forgot password flow.
func (uc *UserUsecase) ForgotPassword(ctx context.Context, email string) error {
	const errUserNotFound = "not found" // Placeholder
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == errUserNotFound {
			uc.logger.Warnf("password reset requested for non-existent email: %s", email)
			return nil
		}
		uc.logger.Errorf("failed to retrieve user for forgot password: %v", err)
		return errors.New("internal server error")
	}

	// Generate a password reset token/link
	resetToken, err := uc.jwtService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		uc.logger.Errorf("failed to generate password reset token: %v", err)
		return errors.New("failed to initiate password reset")
	}

	// Hash the token before storing it to match the schema
	hashedResetToken := uc.hasher.HashString(resetToken)

	// Store the token
	tokenEntity := &entity.Token{
		ID:        uc.uuidGenerator.NewUUID(),
		UserID:    user.ID,
		TokenType: entity.TokenTypePasswordReset,
		TokenHash: hashedResetToken,
		ExpiresAt: time.Now().Add(uc.cfg.GetPasswordResetTokenExpiry()),
		CreatedAt: time.Now(),
		Revoke:    false,
	}
	if err := uc.tokenRepo.CreateToken(ctx, tokenEntity); err != nil {
		uc.logger.Errorf("failed to store password reset token for user %s: %v", user.ID, err)
		return errors.New("failed to initiate password reset")
	}

	// The reset link should use the unhashed token
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.cfg.GetAppBaseURL(), resetToken)
	if err := uc.mailService.SendPasswordResetEmail(user.Email, user.Username, resetLink); err != nil {
		uc.logger.Errorf("failed to send password reset email to %s: %v", user.Email, err)
		return errors.New("failed to send password reset email")
	}

	return nil
}

// ResetPassword handles the password reset flow using a password reset token.
func (uc *UserUsecase) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	// Parse the reset token to get the user claims.
	claims, err := uc.jwtService.ParsePasswordResetToken(resetToken)
	if err != nil {
		return fmt.Errorf("invalid or expired password reset token: %w", err)
	}

	const errTokenNotFound = "not found" // Placeholder
	// Retrieve the stored token using the UserID from the claims.
	storedToken, err := uc.tokenRepo.GetTokenByUserID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errTokenNotFound {
			return errors.New("password reset token not found or invalidated")
		}
		uc.logger.Errorf("failed to retrieve stored reset token: %v", err)
		return errors.New("internal server error")
	}

	// Check if the token has been revoked.
	if storedToken.Revoke {
		return errors.New("password reset token has been used or revoked")
	}

	// Check if the token has expired.
	if storedToken.ExpiresAt.Before(time.Now()) {
		// Delete the expired token to prevent future use.
		_ = uc.tokenRepo.DeleteToken(ctx, storedToken.ID)
		return errors.New("password reset token has expired")
	}

	// Validate the provided reset token against the stored hash.
	if !uc.hasher.CheckHash(resetToken, storedToken.TokenHash) {
		uc.logger.Warnf("password reset token hash mismatch for user %s", claims.UserID)
		// Delete the stored token to prevent further attempts with an invalid token.
		_ = uc.tokenRepo.DeleteToken(ctx, storedToken.ID)
		return errors.New("invalid password reset token")
	}

	// Hash the new password before updating the user.
	hashedNewPassword, err := uc.hasher.HashPassword(newPassword)
	if err != nil {
		uc.logger.Errorf("failed to hash new password: %v", err)
		return errors.New("failed to update password")
	}

	// Update the user's password.
	if err := uc.userRepo.UpdateUserPassword(ctx, claims.UserID, hashedNewPassword); err != nil {
		uc.logger.Errorf("failed to update password for user %s: %v", claims.UserID, err)
		return errors.New("failed to update password")
	}

	// Invalidate the password reset token by deleting it.
	if err := uc.tokenRepo.DeleteToken(ctx, storedToken.ID); err != nil {
		uc.logger.Errorf("failed to delete used password reset token for user %s: %v", claims.UserID, err)
	}

	// Return success, confirming the change.
	return nil
}

// VerifyEmail handles the email verification process.
func (uc *UserUsecase) VerifyEmail(ctx context.Context, token string) error {
	// Parse the email verification token to get the user claims.
	claims, err := uc.jwtService.ParseEmailVerificationToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired email verification token: %w", err)
	}

	const errTokenNotFound = "not found" // Placeholder
	// Retrieve the stored email verification token using the UserID from the claims.
	storedEmailToken, err := uc.emailVerificationTokenRepo.GetEmailVerificationTokenByUserID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errTokenNotFound {
			return errors.New("email verification token not found or invalidated")
		}
		uc.logger.Errorf("failed to retrieve stored email verification token: %v", err)
		return errors.New("internal server error")
	}

	// Check if the token has already been used.
	if storedEmailToken.Used {
		return errors.New("email verification token has already been used")
	}

	// Check if the token has expired.
	if storedEmailToken.ExpiresAt.Before(time.Now()) {
		// Mark the token as used/expired to prevent future use.
		_ = uc.emailVerificationTokenRepo.UpdateEmailVerificationTokenUsedStatus(ctx, storedEmailToken.ID, true)
		return errors.New("email verification token has expired")
	}

	// Validate the provided token against the stored hash.
	if !uc.hasher.CheckHash(token, storedEmailToken.TokenHash) {
		uc.logger.Warnf("email verification token hash mismatch for user %s", claims.UserID)
		// Mark the token as used/invalid to prevent further attempts.
		_ = uc.emailVerificationTokenRepo.UpdateEmailVerificationTokenUsedStatus(ctx, storedEmailToken.ID, true)
		return errors.New("invalid email verification token")
	}

	// Retrieve the user.
	const errUserNotFound = "not found" // Placeholder
	user, err := uc.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errUserNotFound {
			return errors.New("user not found for verification")
		}
		uc.logger.Errorf("failed to retrieve user for email verification: %v", err)
		return errors.New("internal server error")
	}

	// Check if the user is already active.
	if user.IsActive {
		// Mark the token as used even if user is already active.
		_ = uc.emailVerificationTokenRepo.UpdateEmailVerificationTokenUsedStatus(ctx, storedEmailToken.ID, true)
		return errors.New("email already verified")
	}

	// Activate the user's account.
	user.IsActive = true
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		uc.logger.Errorf("failed to activate user %s: %v", user.ID, err)
		return errors.New("failed to activate account")
	}

	// Mark the email verification token as used after successful verification.
	if err := uc.emailVerificationTokenRepo.UpdateEmailVerificationTokenUsedStatus(ctx, storedEmailToken.ID, true); err != nil {
		uc.logger.Errorf("failed to mark email verification token %s as used: %v", storedEmailToken.ID, err)
	}

	return nil
}

// Logout handles user logout.
func (uc *UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	// Parse the refresh token to get the user claims, which gives us the UserID.
	claims, err := uc.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		uc.logger.Warnf("failed to parse refresh token on logout, assuming it's already invalid: %v", err)
		return nil
	}

	const errTokenNotFound = "not found" // Placeholder

	// Retrieve the stored token by UserID to get its database ID.
	storedToken, err := uc.tokenRepo.GetTokenByUserID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == errTokenNotFound {
			uc.logger.Warnf("refresh token for user %s not found during logout, assuming it's already deleted", claims.UserID)
			return nil
		}
		uc.logger.Errorf("failed to retrieve stored refresh token for user %s: %v", claims.UserID, err)
		return errors.New("internal server error")
	}

	// Delete the token from the database.
	if err := uc.tokenRepo.DeleteToken(ctx, storedToken.ID); err != nil {
		uc.logger.Errorf("failed to delete refresh token for user %s: %v", claims.UserID, err)
		return errors.New("failed to delete token")
	}

	return nil
}

// PromoteUser promotes a user to an Admin role.
func (uc *UserUsecase) PromoteUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	const errUserNotFound = "not found" // Placeholder
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err.Error() == errUserNotFound {
			return nil, errors.New("user not found")
		}
		uc.logger.Errorf("failed to retrieve user for promotion: %v", err)
		return nil, errors.New("internal server error")
	}

	if user.Role == entity.UserRoleAdmin {
		return user, errors.New("user is already an admin")
	}

	user.Role = entity.UserRoleAdmin
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		uc.logger.Errorf("failed to promote user %s: %v", userID, err)
		return nil, errors.New("failed to promote user")
	}

	return user, nil
}

// DemoteUser demotes an Admin back to a regular user (member).
func (uc *UserUsecase) DemoteUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	const errUserNotFound = "not found" // Placeholder
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err.Error() == errUserNotFound {
			return nil, errors.New("user not found")
		}
		uc.logger.Errorf("failed to retrieve user for demotion: %v", err)
		return nil, errors.New("internal server error")
	}

	if user.Role == entity.UserRoleUser {
		return user, errors.New("user is already a regular member")
	}

	user.Role = entity.UserRoleUser
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		uc.logger.Errorf("failed to demote user %s: %v", userID, err)
		return nil, errors.New("failed to demote user")
	}

	return user, nil
}

// UpdateProfile allows a registered user to update their profile details.
func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*entity.User, error) {
	const errUserNotFound = "not found" // Placeholder
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err.Error() == errUserNotFound {
			return nil, errors.New("user not found")
		}
		uc.logger.Errorf("failed to retrieve user for profile update: %v", err)
		return nil, errors.New("internal server error")
	}

	// Apply updates based on the map.
	if val, ok := updates["username"]; ok {
		if username, isString := val.(string); isString {
			const errUserNotFound = "not found" // Placeholder
			existingUserByUsername, err := uc.userRepo.GetUserByUsername(ctx, username)
			if err != nil && err.Error() != errUserNotFound {
				uc.logger.Errorf("failed to check for existing username during update: %v", err)
				return nil, errors.New("internal server error")
			}
			if existingUserByUsername != nil && existingUserByUsername.ID != userID {
				return nil, fmt.Errorf("username %s already taken", username)
			}
			user.Username = username
		}
	}
	if val, ok := updates["email"]; ok {
		if email, isString := val.(string); isString {
			if err := uc.validator.ValidateEmail(email); err != nil {
				return nil, fmt.Errorf("invalid email format: %w", err)
			}
			const errUserNotFound = "not found" // Placeholder
			existingUserByEmail, err := uc.userRepo.GetUserByEmail(ctx, email)
			if err != nil && err.Error() != errUserNotFound {
				uc.logger.Errorf("failed to check for existing email during update: %v", err)
				return nil, errors.New("internal server error")
			}
			if existingUserByEmail != nil && existingUserByEmail.ID != userID {
				return nil, fmt.Errorf("email %s already taken", email)
			}
			user.Email = email
		}
	}
	if val, ok := updates["firstName"]; ok {
		if firstName, isString := val.(string); isString {
			if firstName == "" {
				user.FirstName = nil
			} else {
				user.FirstName = &firstName
			}
		}
	}
	if val, ok := updates["lastName"]; ok {
		if lastName, isString := val.(string); isString {
			if lastName == "" {
				user.LastName = nil
			} else {
				user.LastName = &lastName
			}
		}
	}
	if val, ok := updates["avatarURL"]; ok {
		if avatarURL, isString := val.(string); isString {
			if avatarURL == "" {
				user.AvatarURL = nil
			} else {
				user.AvatarURL = &avatarURL
			}
		}
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		uc.logger.Errorf("failed to update profile for user %s: %v", userID, err)
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}

// login with OAuth2

func (uc *UserUsecase) LoginWithOAuth(ctx context.Context, fName, lName, email string) (string, string, error) {
	existingUser, err := uc.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		var generatedUsername string
		if fName != "" && lName != "" {
			generatedUsername = fName + "." + lName
		} else if fName != "" {
			generatedUsername = fName
		} else if lName != "" {
			generatedUsername = lName
		} else {
			generatedUsername = strings.Split(email, "@")[0]
		}
		newUser := &entity.User{
			ID:        uuid.New(),
			Username:  generatedUsername,
			Email:     email,
			Role:      entity.DefaultRole(),
			IsActive:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FirstName: &fName,
			LastName:  &lName,
		}

		err := uc.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			return "", "", fmt.Errorf("failed to create new user: %w", err)
		}
		existingUser = newUser
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(existingUser.ID, existingUser.Role)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenID := uuid.New()
	refreshToken, err := uc.jwtService.GenerateRefreshToken(existingUser.ID, existingUser.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokType, err := entity.SetTokenType("refresh")

	if err != nil {
		return "", "", fmt.Errorf("invaild token type assigned to the refresh token: %w", err)
	}

	refreshTokenRecord := &entity.Token{
		ID:        refreshTokenID,
		UserID:    existingUser.ID,
		TokenType: tokType,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().UTC().Add(168 * time.Hour),
		CreatedAt: time.Now().UTC(),
		Revoke:    false,
	}
	err = uc.tokenRepo.CreateToken(ctx, refreshTokenRecord)
	if err != nil {
		return "", "", fmt.Errorf("failed to save refresh token record: %w", err)
	}
	return accessToken, refreshToken, nil
}
