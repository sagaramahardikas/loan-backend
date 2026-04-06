package entity

type GetResponse struct {
	User User `json:"user"`
}

//go:generate enumer -type=UserStatus -trimprefix=UserStatus -transform=kebab
type UserStatus int8

const (
	UserStatusUnspecified UserStatus = iota
	UserStatusInactive
	UserStatusActive
	UserStatusDelinquent
)

type User struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Status   UserStatus `json:"status"`
}
