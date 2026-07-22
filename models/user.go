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
	RoleCustomer     UserRole = "customer"
	RoleMerchant     UserRole = "merchant"
	RoleAdmin        UserRole = "admin"
	UserPerPageLimit int      = 20
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

func GetUsers(search, role string, offset int) (*[]User, error) {
	query := `SELECT
		id,
		name,
		email,
		role,
		created_at
	FROM users`

	var args []any

	if search != "" {
		query += ` WHERE name LIKE ? OR email LIKE ?`
		args = append(args, "%"+search+"%")
		args = append(args, "%"+search+"%")

		if role != "" {
			query += ` AND role LIKE ?`
			args = append(args, role)
		}
	} else if role != "" {
		query += ` WHERE role LIKE ?`
		args = append(args, role)
	}

	query += ` LIMIT ? OFFSET ?`
	args = append(args, UserPerPageLimit)
	args = append(args, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func UpdateRole(id int64, role string) error {
	query := `UPDATE users SET
		role = ?
	WHERE id = ?`
	_, err := database.DB.Exec(query, role, id)

	return err
}
