package errs

// User-friendly error messages
const (
	MsgInvalidRequest         = "Invalid request. Please check your input."
	MsgInternalError          = "An unexpected error occurred. Please try again later."
	MsgInvalidCredentials     = "Invalid email or password."
	MsgEmailAlreadyRegistered = "This email is already registered."
	MsgUserNotFound           = "User not found."
	MsgTokenExpired           = "Your session has expired. Please login again."
	MsgInvalidToken           = "Invalid authentication token."
	MsgUnauthorized           = "You are not authorized to perform this action."
	MsgSessionNotFound        = "Session not found. Please login again."
	MsgAddressNotFound        = "Address not found."
	MsgRefreshFailed          = "Failed to refresh token. Please login again."
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Error: message}
}