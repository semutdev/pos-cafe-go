package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// User represents the users table.
type User struct {
	ID                        int     `json:"id" db:"id"`
	IPAddress                 string  `json:"ip_address" db:"ip_address"`
	Username                  string  `json:"username" db:"username"`
	Password                  string  `json:"password" db:"password"`
	Email                     string  `json:"email" db:"email"`
	ActivationSelector        *string `json:"activation_selector" db:"activation_selector"`
	ActivationCode            *string `json:"activation_code" db:"activation_code"`
	ForgottenPasswordSelector *string `json:"forgotten_password_selector" db:"forgotten_password_selector"`
	ForgottenPasswordCode     *string `json:"forgotten_password_code" db:"forgotten_password_code"`
	ForgottenPasswordTime     *int    `json:"forgotten_password_time" db:"forgotten_password_time"`
	RememberSelector          *string `json:"remember_selector" db:"remember_selector"`
	RememberCode              *string `json:"remember_code" db:"remember_code"`
	CreatedOn                 int     `json:"created_on" db:"created_on"`
	LastLogin                 *int    `json:"last_login" db:"last_login"`
	Active                    *bool   `json:"active" db:"active"`
	FirstName                 *string `json:"first_name" db:"first_name"`
	LastName                  *string `json:"last_name" db:"last_name"`
	Company                   *string `json:"company" db:"company"`
	Phone                     *string `json:"phone" db:"phone"`
}

// TableName returns the database table name for User.
func (User) TableName() string {
	return "users"
}

// UserModel provides data access methods for users.
type UserModel struct{}

// NewUserModel creates a new UserModel instance.
func NewUserModel() *UserModel {
	return &UserModel{}
}

// Find retrieves a single User by primary key.
func (m *UserModel) Find(id int) (*User, error) {
	var item User
	err := database.Table("users").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all User records.
func (m *UserModel) FindAll() ([]User, error) {
	var items []User
	err := database.Table("users").Get(&items)
	return items, err
}

// Insert inserts a new User record into the database.
func (m *UserModel) Insert(user *User) error {
	return database.Table("users").Insert(map[string]any{
		"id":                          user.ID,
		"ip_address":                  user.IPAddress,
		"username":                    user.Username,
		"password":                    user.Password,
		"email":                       user.Email,
		"activation_selector":         user.ActivationSelector,
		"activation_code":             user.ActivationCode,
		"forgotten_password_selector": user.ForgottenPasswordSelector,
		"forgotten_password_code":     user.ForgottenPasswordCode,
		"forgotten_password_time":     user.ForgottenPasswordTime,
		"remember_selector":           user.RememberSelector,
		"remember_code":               user.RememberCode,
		"created_on":                  user.CreatedOn,
		"last_login":                  user.LastLogin,
		"active":                      user.Active,
		"first_name":                  user.FirstName,
		"last_name":                   user.LastName,
		"company":                     user.Company,
		"phone":                       user.Phone,
	})
}

// Update updates an existing User record in the database.
func (m *UserModel) Update(user *User) error {
	return database.Table("users").Where("id", user.ID).Update(map[string]any{
		"ip_address":                  user.IPAddress,
		"username":                    user.Username,
		"password":                    user.Password,
		"email":                       user.Email,
		"activation_selector":         user.ActivationSelector,
		"activation_code":             user.ActivationCode,
		"forgotten_password_selector": user.ForgottenPasswordSelector,
		"forgotten_password_code":     user.ForgottenPasswordCode,
		"forgotten_password_time":     user.ForgottenPasswordTime,
		"remember_selector":           user.RememberSelector,
		"remember_code":               user.RememberCode,
		"created_on":                  user.CreatedOn,
		"last_login":                  user.LastLogin,
		"active":                      user.Active,
		"first_name":                  user.FirstName,
		"last_name":                   user.LastName,
		"company":                     user.Company,
		"phone":                       user.Phone,
	})
}

// Delete deletes a User record by primary key.
func (m *UserModel) Delete(id int) error {
	return database.Table("users").Where("id", id).Delete()
}
