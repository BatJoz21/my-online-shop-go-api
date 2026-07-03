package models

type StoredUserData struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}
