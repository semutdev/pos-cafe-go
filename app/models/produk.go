package models

import (
	"time"

	"github.com/semutdev/goigniter/system/libraries/database"
)

// Produk represents the produk table.
type Produk struct {
	ID               int        `json:"id" db:"id"`
	IDProdukKategori int        `json:"id_produk_kategori" db:"id_produk_kategori"`
	CodeProduk       string     `json:"code_produk" db:"code_produk"`
	NamaProduk       string     `json:"nama_produk" db:"nama_produk"`
	Deskripsi        *string    `json:"deskripsi" db:"deskripsi"`
	Harga            float64    `json:"harga" db:"harga"`
	Stok             int        `json:"stok" db:"stok"`
	Gambar           *string    `json:"gambar" db:"gambar"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	IsDeleted        bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// TableName returns the database table name for Produk.
func (Produk) TableName() string {
	return "produk"
}

// ProdukModel provides data access methods for produk.
type ProdukModel struct{}

// NewProdukModel creates a new ProdukModel instance.
func NewProdukModel() *ProdukModel {
	return &ProdukModel{}
}

// Find retrieves a single Produk by primary key.
func (m *ProdukModel) Find(id int) (*Produk, error) {
	var item Produk
	err := database.Table("produk").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all Produk records.
func (m *ProdukModel) FindAll() ([]Produk, error) {
	var items []Produk
	err := database.Table("produk").Get(&items)
	return items, err
}

// Insert inserts a new Produk record into the database.
func (m *ProdukModel) Insert(produk *Produk) error {
	return database.Table("produk").Insert(map[string]any{
		"id":                 produk.ID,
		"id_produk_kategori": produk.IDProdukKategori,
		"code_produk":        produk.CodeProduk,
		"nama_produk":        produk.NamaProduk,
		"deskripsi":          produk.Deskripsi,
		"harga":              produk.Harga,
		"stok":               produk.Stok,
		"gambar":             produk.Gambar,
		"is_active":          produk.IsActive,
		"is_deleted":         produk.IsDeleted,
		"deleted_at":         produk.DeletedAt,
		"updated_at":         produk.UpdatedAt,
		"created_at":         produk.CreatedAt,
	})
}

// Update updates an existing Produk record in the database.
func (m *ProdukModel) Update(produk *Produk) error {
	return database.Table("produk").Where("id", produk.ID).Update(map[string]any{
		"id_produk_kategori": produk.IDProdukKategori,
		"code_produk":        produk.CodeProduk,
		"nama_produk":        produk.NamaProduk,
		"deskripsi":          produk.Deskripsi,
		"harga":              produk.Harga,
		"stok":               produk.Stok,
		"gambar":             produk.Gambar,
		"is_active":          produk.IsActive,
		"is_deleted":         produk.IsDeleted,
		"deleted_at":         produk.DeletedAt,
		"updated_at":         produk.UpdatedAt,
		"created_at":         produk.CreatedAt,
	})
}

// Delete deletes a Produk record by primary key.
func (m *ProdukModel) Delete(id int) error {
	return database.Table("produk").Where("id", id).Delete()
}
