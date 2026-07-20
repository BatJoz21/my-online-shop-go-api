package models

import (
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/database"
)

type User struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         UserRole   `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleMerchant UserRole = "merchant"
	RoleAdmin    UserRole = "admin"
)

func (u *User) Save() error {
	query := `INSERT INTO users(name, email, password_hash, role) VALUES (?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Name, u.Email, u.PasswordHash, u.Role)
	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetUserDataForSession(email string) (*StoredUserData, error) {
	query := `SELECT id, name, email, role FROM users WHERE email = ?`
	row := database.DB.QueryRow(query, email)

	var userData StoredUserData
	err := row.Scan(&userData.ID, &userData.Name, &userData.Email, &userData.Role)

	return &userData, err
}

func GetUserDataForRefreshToken(id int64) (*StoredUserData, error) {
	query := `SELECT id, name, email, role FROM users WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var userData StoredUserData
	err := row.Scan(&userData.ID, &userData.Name, &userData.Email, &userData.Role)

	return &userData, err
}

func GetUsers() (*[]User, error) {
	query := `SELECT
		id,
		name,
		email,
		role,
		created_at
	FROM users`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUser(id int64) (*User, error) {
	query := `SELECT
		id,
		name,
		email,
		role,
		created_at
	FROM users WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func GetUserForProfile(id int64) (*User, error) {
	query := `SELECT
		name,
		email,
	FROM users WHERE id = ?`
	row := database.DB.QueryRow(query)

	var u User
	if err := row.Scan(&u.Name, &u.Email); err != nil {
		return nil, err
	}

	return &u, nil
}

func UpdateRole(id int64, role string) error {
	query := `UPDATE users SET
		role = ?
	WHERE id = ?`
	_, err := database.DB.Exec(query, role, id)

	return err
}
