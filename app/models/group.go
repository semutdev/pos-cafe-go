package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// Group represents the groups table.
type Group struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// TableName returns the database table name for Group.
func (Group) TableName() string {
	return "groups"
}

// GroupModel provides data access methods for groups.
type GroupModel struct{}

// NewGroupModel creates a new GroupModel instance.
func NewGroupModel() *GroupModel {
	return &GroupModel{}
}

// Find retrieves a single Group by primary key.
func (m *GroupModel) Find(id int) (*Group, error) {
	var item Group
	err := database.Table("groups").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all Group records.
func (m *GroupModel) FindAll() ([]Group, error) {
	var items []Group
	err := database.Table("groups").Get(&items)
	return items, err
}

// Insert inserts a new Group record into the database.
func (m *GroupModel) Insert(group *Group) error {
	return database.Table("groups").Insert(map[string]any{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
	})
}

// Update updates an existing Group record in the database.
func (m *GroupModel) Update(group *Group) error {
	return database.Table("groups").Where("id", group.ID).Update(map[string]any{
		"name":        group.Name,
		"description": group.Description,
	})
}

// Delete deletes a Group record by primary key.
func (m *GroupModel) Delete(id int) error {
	return database.Table("groups").Where("id", id).Delete()
}
