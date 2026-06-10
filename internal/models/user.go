package models

// CreateUserRequest is the request body for creating a new user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the request body for updating an existing user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UserResponse is the response body for create and update operations (no age).
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserWithAgeResponse is the response body for get and list operations (includes age).
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginationMeta holds pagination metadata returned alongside list results.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// ListUsersResponse wraps the user list with pagination metadata.
type ListUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

// ErrorResponse is a standard JSON error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}
