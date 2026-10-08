package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// UsersGroup represents the users_groups table.
type UsersGroup struct {
	ID      int `json:"id" db:"id"`
	UserID  int `json:"user_id" db:"user_id"`
	GroupID int `json:"group_id" db:"group_id"`
}

// TableName returns the database table name for UsersGroup.
func (UsersGroup) TableName() string {
	return "users_groups"
}

// UsersGroupModel provides data access methods for users_groups.
type UsersGroupModel struct{}

// NewUsersGroupModel creates a new UsersGroupModel instance.
func NewUsersGroupModel() *UsersGroupModel {
	return &UsersGroupModel{}
}

// Find retrieves a single UsersGroup by primary key.
func (m *UsersGroupModel) Find(id int) (*UsersGroup, error) {
	var item UsersGroup
	err := database.Table("users_groups").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all UsersGroup records.
func (m *UsersGroupModel) FindAll() ([]UsersGroup, error) {
	var items []UsersGroup
	err := database.Table("users_groups").Get(&items)
	return items, err
}

// Insert inserts a new UsersGroup record into the database.
func (m *UsersGroupModel) Insert(usersGroup *UsersGroup) error {
	return database.Table("users_groups").Insert(map[string]any{
		"id":       usersGroup.ID,
		"user_id":  usersGroup.UserID,
		"group_id": usersGroup.GroupID,
	})
}

// Update updates an existing UsersGroup record in the database.
func (m *UsersGroupModel) Update(usersGroup *UsersGroup) error {
	return database.Table("users_groups").Where("id", usersGroup.ID).Update(map[string]any{
		"user_id":  usersGroup.UserID,
		"group_id": usersGroup.GroupID,
	})
}

// Delete deletes a UsersGroup record by primary key.
func (m *UsersGroupModel) Delete(id int) error {
	return database.Table("users_groups").Where("id", id).Delete()
}
