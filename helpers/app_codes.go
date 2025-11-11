package helpers

// Reference: docs/app_codes.md

const (
	Success = 1 + iota
	SRequestAccepted
	SResourceCreatedSuccessfully
	SResourceUpdatedSuccessfully
	SResourceDeletedSuccessfully
)

const (
	Ignored          = 0
	UnspecifiedError = -1
)

const (
	EInvalidRequest = -2 - iota

	EZaloSignatureMismatch = -90
)

const (
	ESignatureInvalid = -100 - iota
	EMissingAuthorizationHeader
	EInvalidAccessToken
	EExpiredAccessToken
	EInvalidRefreshToken
	EAccessDenied
	EAccountSuspended
	EJWTGenerationFailed
)

const (
	EUserNotFound = -200 - iota
	EUserAlreadyExists
	EUsernameTaken
	ETooManyLoginAttempts
)

const (
	ENotEnoughCoins = -300 - iota
	ENotEnoughGems
	EEnergyDepleted
)

const (
	EResourceNotFound = -400 - iota
	EResourceAlreadyExists
	EDatabaseError
	EDataIntegrityViolation
	ERecordNotAllowedToModify
)

const (
	EExternalServiceError = -500 - iota
	EExternalServiceTimeout
	EThirdPartyServiceUnavailable
)

var Messages = map[int]string{ //nolint:gochecknoglobals // global map of code-message pairs
	Success:                      "Success",
	SRequestAccepted:             "Request accepted",
	SResourceCreatedSuccessfully: "Resource created successfully",
	SResourceUpdatedSuccessfully: "Resource updated successfully",
	SResourceDeletedSuccessfully: "Resource deleted successfully",

	EInvalidRequest: "Invalid request",

	EZaloSignatureMismatch: "Zalo signature mismatch",

	ESignatureInvalid:           "Invalid signature",
	EMissingAuthorizationHeader: "Missing authorization header",
	EInvalidAccessToken:         "Invalid access token",
	EExpiredAccessToken:         "Expired access token",
	EInvalidRefreshToken:        "Invalid refresh token",
	EAccessDenied:               "Access denied",
	EAccountSuspended:           "Account suspended",

	EUserNotFound:         "User not found",
	EUserAlreadyExists:    "User already exists",
	EUsernameTaken:        "Username is already taken",
	ETooManyLoginAttempts: "Too many login attempts",

	ENotEnoughCoins: "Not enough coins",
	ENotEnoughGems:  "Not enough gems",
	EEnergyDepleted: "Energy depleted",

	EResourceNotFound:         "Resource not found",
	EResourceAlreadyExists:    "Resource already exists",
	EDatabaseError:            "Database error",
	EDataIntegrityViolation:   "Data integrity violation",
	ERecordNotAllowedToModify: "Record not allowed to modify",

	EExternalServiceError:         "External service error",
	EExternalServiceTimeout:       "External service timeout",
	EThirdPartyServiceUnavailable: "Third-party service unavailable",
}

func Message(code int) string {
	if msg, ok := Messages[code]; ok {
		return msg
	}
	return "Unknown code"
}
