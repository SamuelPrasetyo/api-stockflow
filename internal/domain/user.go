package domain

import "time"

// UserRole represents user role type
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleStaff   UserRole = "staff"
)

// User represents user entity
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	Fullname  string    `json:"fullname"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CanApprove checks if user can approve a purchase based on amount
func (u *User) CanApprove(amount float64) bool {
	// Admin and Manager can approve any amount
	if u.Role == RoleAdmin || u.Role == RoleManager {
		return true
	}
	// Staff can only approve if amount is less than or equal to 10,000,000
	if u.Role == RoleStaff && amount <= 10000000 {
		return true
	}
	return false
}

// HasFullAccess checks if user has full access to all menus
func (u *User) HasFullAccess() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}
