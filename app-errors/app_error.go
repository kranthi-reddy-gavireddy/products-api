package apperrors

import "net/http"

type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	HTTPStatus int         `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func ValidationError(details interface{}) *AppError {
	return &AppError{
		Code:       ValidationErrorCode,
		Message:    VALIDATION_ERROR_MSG,
		Details:    details,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NotFoundError(details interface{}) *AppError {
	return &AppError{
		Code:       NotFoundErrorCode,
		Message:    NOT_FOUND_ERROR_MSG,
		Details:    details,
		HTTPStatus: http.StatusNotFound,
	}
}

func InternalServerError(details interface{}) *AppError {
	return &AppError{
		Code:       InternalErrorCode,
		Message:    INTERNAL_ERROR_MSG,
		Details:    details,
		HTTPStatus: http.StatusInternalServerError,
	}
}

func UnauthorizedError(details interface{}) *AppError {
	return &AppError{
		Code:       UnauthorizedErrorCode,
		Message:    UNAUTHORIZED_ERROR_MSG,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func ForbiddenError(details interface{}) *AppError {
	return &AppError{
		Code:       ForbiddenErrorCode,
		Message:    FORBIDDEN_ERROR_MSG,
		HTTPStatus: http.StatusForbidden,
	}
}

func DatabaseError() *AppError {
	return &AppError{
		Code:       DatabaseErrorCode,
		Message:    DATABASE_ERROR_MSG,
		HTTPStatus: http.StatusInternalServerError,
	}
}

func BadRequestError(details interface{}) *AppError {
	return &AppError{
		Code:       BadRequestErrorCode,
		Message:    BAD_REQUEST_MSG,
		Details:    details,
		HTTPStatus: http.StatusBadRequest,
	}
}
