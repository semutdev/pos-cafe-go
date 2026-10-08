package models

import (
	"time"

	"github.com/semutdev/goigniter/system/libraries/database"
)

// Transaksi represents the transaksi table.
type Transaksi struct {
	ID               int       `json:"id" db:"id"`
	Tanggal          time.Time `json:"tanggal" db:"tanggal"`
	NoInvoice        string    `json:"no_invoice" db:"no_invoice"`
	IDUser           int       `json:"id_user" db:"id_user"`
	TotalHarga       float64   `json:"total_harga" db:"total_harga"`
	Diskon           int       `json:"diskon" db:"diskon"`
	TotalAkhir       float64   `json:"total_akhir" db:"total_akhir"`
	UangBayar        float64   `json:"uang_bayar" db:"uang_bayar"`
	MetodePembayaran string    `json:"metode_pembayaran" db:"metode_pembayaran"`
	TipePesanan      *string   `json:"tipe_pesanan" db:"tipe_pesanan"`
	Status           string    `json:"status" db:"status"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// TableName returns the database table name for Transaksi.
func (Transaksi) TableName() string {
	return "transaksi"
}

// TransaksiModel provides data access methods for transaksi.
type TransaksiModel struct{}

// NewTransaksiModel creates a new TransaksiModel instance.
func NewTransaksiModel() *TransaksiModel {
	return &TransaksiModel{}
}

// Find retrieves a single Transaksi by primary key.
func (m *TransaksiModel) Find(id int) (*Transaksi, error) {
	var item Transaksi
	err := database.Table("transaksi").Where("id", id).First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll retrieves all Transaksi records.
func (m *TransaksiModel) FindAll() ([]Transaksi, error) {
	var items []Transaksi
	err := database.Table("transaksi").Get(&items)
	return items, err
}

// Insert inserts a new Transaksi record into the database.
func (m *TransaksiModel) Insert(transaksi *Transaksi) error {
	return database.Table("transaksi").Insert(map[string]any{
		"id":                transaksi.ID,
		"tanggal":           transaksi.Tanggal,
		"no_invoice":        transaksi.NoInvoice,
		"id_user":           transaksi.IDUser,
		"total_harga":       transaksi.TotalHarga,
		"diskon":            transaksi.Diskon,
		"total_akhir":       transaksi.TotalAkhir,
		"uang_bayar":        transaksi.UangBayar,
		"metode_pembayaran": transaksi.MetodePembayaran,
		"tipe_pesanan":      transaksi.TipePesanan,
		"status":            transaksi.Status,
		"updated_at":        transaksi.UpdatedAt,
		"created_at":        transaksi.CreatedAt,
	})
}

// Update updates an existing Transaksi record in the database.
func (m *TransaksiModel) Update(transaksi *Transaksi) error {
	return database.Table("transaksi").Where("id", transaksi.ID).Update(map[string]any{
		"tanggal":           transaksi.Tanggal,
		"no_invoice":        transaksi.NoInvoice,
		"id_user":           transaksi.IDUser,
		"total_harga":       transaksi.TotalHarga,
		"diskon":            transaksi.Diskon,
		"total_akhir":       transaksi.TotalAkhir,
		"uang_bayar":        transaksi.UangBayar,
		"metode_pembayaran": transaksi.MetodePembayaran,
		"tipe_pesanan":      transaksi.TipePesanan,
		"status":            transaksi.Status,
		"updated_at":        transaksi.UpdatedAt,
		"created_at":        transaksi.CreatedAt,
	})
}

// Delete deletes a Transaksi record by primary key.
func (m *TransaksiModel) Delete(id int) error {
	return database.Table("transaksi").Where("id", id).Delete()
}
