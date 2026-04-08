package client

type UpdateUserRequest struct {
	UserID string `json:"id"`
	Status int    `json:"status"` // 1: inactive, 2: active, 3: delinquent
}

type UpdateUserResponse struct {
	UserID string `json:"id"`
	Status int    `json:"status"` // 1: inactive, 2: active, 3: delinquent
}
