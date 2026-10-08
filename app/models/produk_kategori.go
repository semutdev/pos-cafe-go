package models

import (
	"time"

	"github.com/semutdev/goigniter/system/libraries/database"
)

// ProdukKategori represents the produk_kategori table.
type ProdukKategori struct {
	ID           int       `json:"id" db:"id"`
	CodeKategori string    `json:"code_kategori" db:"code_kategori"`
	NamaKategori string    `json:"nama_kategori" db:"nama_kategori"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the database table name for ProdukKategori.
func (ProdukKategori) TableName() string {
	return "produk_kategori"
}

// ProdukKategoriModel provides data access methods for produk_kategori.
type ProdukKategoriModel struct{}

// NewProdukKategoriModel creates a new ProdukKategoriModel instance.
func NewProdukKategoriModel() *ProdukKategoriModel {
	return &ProdukKategoriModel{}
}

// Find retrieves a single ProdukKategori by primary key.
func (m *ProdukKategoriModel) Find(id int) (*ProdukKategori, error) {
	var item ProdukKategori
	err := database.Table("produk_kategori").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all ProdukKategori records.
func (m *ProdukKategoriModel) FindAll() ([]ProdukKategori, error) {
	var items []ProdukKategori
	err := database.Table("produk_kategori").Get(&items)
	return items, err
}

// Insert inserts a new ProdukKategori record into the database.
func (m *ProdukKategoriModel) Insert(produkKategori *ProdukKategori) error {
	return database.Table("produk_kategori").Insert(map[string]any{
		"id":            produkKategori.ID,
		"code_kategori": produkKategori.CodeKategori,
		"nama_kategori": produkKategori.NamaKategori,
		"is_active":     produkKategori.IsActive,
		"updated_at":    produkKategori.UpdatedAt,
		"created_at":    produkKategori.CreatedAt,
	})
}

// Update updates an existing ProdukKategori record in the database.
func (m *ProdukKategoriModel) Update(produkKategori *ProdukKategori) error {
	return database.Table("produk_kategori").Where("id", produkKategori.ID).Update(map[string]any{
		"code_kategori": produkKategori.CodeKategori,
		"nama_kategori": produkKategori.NamaKategori,
		"is_active":     produkKategori.IsActive,
		"updated_at":    produkKategori.UpdatedAt,
		"created_at":    produkKategori.CreatedAt,
	})
}

// Delete deletes a ProdukKategori record by primary key.
func (m *ProdukKategoriModel) Delete(id int) error {
	return database.Table("produk_kategori").Where("id", id).Delete()
}
