package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// Migration represents the migrations table.
type Migration struct {
	Version int64 `json:"version" db:"version"`
}

// TableName returns the database table name for Migration.
func (Migration) TableName() string {
	return "migrations"
}

// MigrationModel provides data access methods for migrations.
type MigrationModel struct{}

// NewMigrationModel creates a new MigrationModel instance.
func NewMigrationModel() *MigrationModel {
	return &MigrationModel{}
}

// Find retrieves a single Migration by version (table lacks primary key).
func (m *MigrationModel) Find(version int64) (*Migration, error) {
	var item Migration
	err := database.Table("migrations").Where("version", version).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all Migration records.
func (m *MigrationModel) FindAll() ([]Migration, error) {
	var items []Migration
	err := database.Table("migrations").Get(&items)
	return items, err
}

// Insert inserts a new Migration record into the database.
func (m *MigrationModel) Insert(migration *Migration) error {
	return database.Table("migrations").Insert(map[string]any{
		"version": migration.Version,
	})
}

// Update updates an existing Migration record in the database by version (table lacks primary key).
func (m *MigrationModel) Update(migration *Migration) error {
	return database.Table("migrations").Where("version", migration.Version).Update(map[string]any{
		"version": migration.Version,
	})
}

// Delete deletes a Migration record by version (table lacks primary key).
func (m *MigrationModel) Delete(version int64) error {
	return database.Table("migrations").Where("version", version).Delete()
}
