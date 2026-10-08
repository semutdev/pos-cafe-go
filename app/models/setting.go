package models

import (
	"github.com/semutdev/goigniter/system/libraries/database"
)

// Setting represents the setting table.
type Setting struct {
	SettingKey   string  `json:"setting_key" db:"setting_key"`
	SettingValue *string `json:"setting_value" db:"setting_value"`
}

// TableName returns the database table name for Setting.
func (Setting) TableName() string {
	return "setting"
}

// SettingModel provides data access methods for setting.
type SettingModel struct{}

// NewSettingModel creates a new SettingModel instance.
func NewSettingModel() *SettingModel {
	return &SettingModel{}
}

// Find retrieves a single Setting by setting_key (table lacks primary key).
func (m *SettingModel) Find(settingKey string) (*Setting, error) {
	var item Setting
	err := database.Table("setting").Where("setting_key", settingKey).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all Setting records.
func (m *SettingModel) FindAll() ([]Setting, error) {
	var items []Setting
	err := database.Table("setting").Get(&items)
	return items, err
}

// Insert inserts a new Setting record into the database.
func (m *SettingModel) Insert(setting *Setting) error {
	return database.Table("setting").Insert(map[string]any{
		"setting_key":   setting.SettingKey,
		"setting_value": setting.SettingValue,
	})
}

// Update updates an existing Setting record in the database by setting_key (table lacks primary key).
func (m *SettingModel) Update(setting *Setting) error {
	return database.Table("setting").Where("setting_key", setting.SettingKey).Update(map[string]any{
		"setting_value": setting.SettingValue,
	})
}

// Delete deletes a Setting record by setting_key (table lacks primary key).
func (m *SettingModel) Delete(settingKey string) error {
	return database.Table("setting").Where("setting_key", settingKey).Delete()
}
