package datastruct

// Constant Value
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserResponse struct {
	ResponseCode string
	ResponseStatus string
}

type UserDataResponse struct {
	Users []UserRequest
}