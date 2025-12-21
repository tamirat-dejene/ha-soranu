package domain

const (
	InvalidCredentialsMessage     = "Invalid email or secret key"
	RestaurantNotFoundMessage       = "Restaurant not found"
	RestaurantAlreadyExistsMessage  = "Restaurant already exists"
	InvalidRestaurantDataMessage   = "Invalid restaurant data provided"
	RestaurantCreationSuccessMessage = "Restaurant created successfully"
	InvalidSearchDataMessage       = "Invalid search data provided"
	InvalidOrderDataMessage        = "Invalid order data provided"
	OrderNotFoundMessage           = "Order not found"
)

var (
	ErrInvalidCredentials      = NewDomainError(InvalidCredentialsMessage)
	ErrRestaurantNotFound       = NewDomainError(RestaurantNotFoundMessage)
	ErrRestaurantAlreadyExists  = NewDomainError(RestaurantAlreadyExistsMessage)
	ErrInvalidRestaurantData    = NewDomainError(InvalidRestaurantDataMessage)
	ErrRestaurantCreationFailed = NewDomainError("Failed to create restaurant")
	ErrMenuItemNotFound         = NewDomainError("Menu item not found")
	ErrInvalidSearchData       = NewDomainError(InvalidSearchDataMessage)
	ErrInvalidOrderData        = NewDomainError(InvalidOrderDataMessage)
	ErrOrderNotFound           = NewDomainError(OrderNotFoundMessage)
)

type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(message string) error {
	return &DomainError{Message: message}
}