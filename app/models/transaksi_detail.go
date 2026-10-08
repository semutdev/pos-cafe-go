package models

import (
	"time"

	"github.com/semutdev/goigniter/system/libraries/database"
)

// TransaksiDetail represents the transaksi_detail table.
type TransaksiDetail struct {
	ID                 int       `json:"id" db:"id"`
	IDTransaksi        int       `json:"id_transaksi" db:"id_transaksi"`
	IDProduk           int       `json:"id_produk" db:"id_produk"`
	IDProdukKategori   int       `json:"id_produk_kategori" db:"id_produk_kategori"`
	CodeProduk         string    `json:"code_produk" db:"code_produk"`
	NamaProduk         string    `json:"nama_produk" db:"nama_produk"`
	CodeKategori       string    `json:"code_kategori" db:"code_kategori"`
	NamaKategori       string    `json:"nama_kategori" db:"nama_kategori"`
	HargaSaatTransaksi int       `json:"harga_saat_transaksi" db:"harga_saat_transaksi"`
	Jumlah             int       `json:"jumlah" db:"jumlah"`
	Subtotal           int       `json:"subtotal" db:"subtotal"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the database table name for TransaksiDetail.
func (TransaksiDetail) TableName() string {
	return "transaksi_detail"
}

// TransaksiDetailModel provides data access methods for transaksi_detail.
type TransaksiDetailModel struct{}

// NewTransaksiDetailModel creates a new TransaksiDetailModel instance.
func NewTransaksiDetailModel() *TransaksiDetailModel {
	return &TransaksiDetailModel{}
}

// Find retrieves a single TransaksiDetail by primary key.
func (m *TransaksiDetailModel) Find(id int) (*TransaksiDetail, error) {
	var item TransaksiDetail
	err := database.Table("transaksi_detail").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all TransaksiDetail records.
func (m *TransaksiDetailModel) FindAll() ([]TransaksiDetail, error) {
	var items []TransaksiDetail
	err := database.Table("transaksi_detail").Get(&items)
	return items, err
}

// Insert inserts a new TransaksiDetail record into the database.
func (m *TransaksiDetailModel) Insert(transaksiDetail *TransaksiDetail) error {
	return database.Table("transaksi_detail").Insert(map[string]any{
		"id":                   transaksiDetail.ID,
		"id_transaksi":         transaksiDetail.IDTransaksi,
		"id_produk":            transaksiDetail.IDProduk,
		"id_produk_kategori":   transaksiDetail.IDProdukKategori,
		"code_produk":          transaksiDetail.CodeProduk,
		"nama_produk":          transaksiDetail.NamaProduk,
		"code_kategori":        transaksiDetail.CodeKategori,
		"nama_kategori":        transaksiDetail.NamaKategori,
		"harga_saat_transaksi": transaksiDetail.HargaSaatTransaksi,
		"jumlah":               transaksiDetail.Jumlah,
		"subtotal":             transaksiDetail.Subtotal,
		"updated_at":           transaksiDetail.UpdatedAt,
		"created_at":           transaksiDetail.CreatedAt,
	})
}

// Update updates an existing TransaksiDetail record in the database.
func (m *TransaksiDetailModel) Update(transaksiDetail *TransaksiDetail) error {
	return database.Table("transaksi_detail").Where("id", transaksiDetail.ID).Update(map[string]any{
		"id_transaksi":         transaksiDetail.IDTransaksi,
		"id_produk":            transaksiDetail.IDProduk,
		"id_produk_kategori":   transaksiDetail.IDProdukKategori,
		"code_produk":          transaksiDetail.CodeProduk,
		"nama_produk":          transaksiDetail.NamaProduk,
		"code_kategori":        transaksiDetail.CodeKategori,
		"nama_kategori":        transaksiDetail.NamaKategori,
		"harga_saat_transaksi": transaksiDetail.HargaSaatTransaksi,
		"jumlah":               transaksiDetail.Jumlah,
		"subtotal":             transaksiDetail.Subtotal,
		"updated_at":           transaksiDetail.UpdatedAt,
		"created_at":           transaksiDetail.CreatedAt,
	})
}

// Delete deletes a TransaksiDetail record by primary key.
func (m *TransaksiDetailModel) Delete(id int) error {
	return database.Table("transaksi_detail").Where("id", id).Delete()
}
