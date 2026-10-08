package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// LoginAttempt represents the login_attempts table.
type LoginAttempt struct {
	ID        int     `json:"id" db:"id"`
	IPAddress string  `json:"ip_address" db:"ip_address"`
	Login     *string `json:"login" db:"login"`
	Time      *int    `json:"time" db:"time"`
}

// TableName returns the database table name for LoginAttempt.
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// LoginAttemptModel provides data access methods for login_attempts.
type LoginAttemptModel struct{}

// NewLoginAttemptModel creates a new LoginAttemptModel instance.
func NewLoginAttemptModel() *LoginAttemptModel {
	return &LoginAttemptModel{}
}

// Find retrieves a single LoginAttempt by primary key.
func (m *LoginAttemptModel) Find(id int) (*LoginAttempt, error) {
	var item LoginAttempt
	err := database.Table("login_attempts").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all LoginAttempt records.
func (m *LoginAttemptModel) FindAll() ([]LoginAttempt, error) {
	var items []LoginAttempt
	err := database.Table("login_attempts").Get(&items)
	return items, err
}

// Insert inserts a new LoginAttempt record into the database.
func (m *LoginAttemptModel) Insert(loginAttempt *LoginAttempt) error {
	return database.Table("login_attempts").Insert(map[string]any{
		"id":         loginAttempt.ID,
		"ip_address": loginAttempt.IPAddress,
		"login":      loginAttempt.Login,
		"time":       loginAttempt.Time,
	})
}

// Update updates an existing LoginAttempt record in the database.
func (m *LoginAttemptModel) Update(loginAttempt *LoginAttempt) error {
	return database.Table("login_attempts").Where("id", loginAttempt.ID).Update(map[string]any{
		"ip_address": loginAttempt.IPAddress,
		"login":      loginAttempt.Login,
		"time":       loginAttempt.Time,
	})
}

// Delete deletes a LoginAttempt record by primary key.
func (m *LoginAttemptModel) Delete(id int) error {
	return database.Table("login_attempts").Where("id", id).Delete()
}
