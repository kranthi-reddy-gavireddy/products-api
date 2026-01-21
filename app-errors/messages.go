package apperrors

const (
	BadRequestErrorCode    = "BAD_REQUEST"
	ValidationErrorCode    = "VALIDATION_ERROR"
	NotFoundErrorCode      = "NOT_FOUND"
	InternalErrorCode      = "INTERNAL_ERROR"
	UnauthorizedErrorCode  = "UNAUTHORIZED"
	ForbiddenErrorCode     = "FORBIDDEN"
	DatabaseErrorCode      = "DATABASE_ERROR"
	BAD_REQUEST_MSG        = "Bad request"
	VALIDATION_ERROR_MSG   = "Validation failed"
	NOT_FOUND_ERROR_MSG    = "Resource not found"
	INTERNAL_ERROR_MSG     = "Internal server error"
	UNAUTHORIZED_ERROR_MSG = "Unauthorized access"
	FORBIDDEN_ERROR_MSG    = "Forbidden access"
	DATABASE_ERROR_MSG     = "Database error"
)

const (
	PRODUCT_NOT_FOUND     = "Product not found with the given ID %s"
	INVALID_PARAMS        = "%s Cannot be an empty value"
	INSUFFICIENT_QUANTITY = "Insufficient product quantity available for product ID %s"
)
