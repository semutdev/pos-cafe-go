package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/libraries/database"
	"pos-cafe-go/app/models"
)

type ProdukController struct {
	core.Controller
}

func (ctrl *ProdukController) Routes() map[string]string {
	return map[string]string{
		"GetData":          "get_data_produk",
		"Tambah":           "tambah",
		"Simpan":           "simpan",
		"Edit":             "edit/:id",
		"Update":           "update/:id",
		"Detail":           "detail/:id",
		"Hapus":            "hapus/:id",
		"BulkDelete":       "bulk_delete",
		"DownloadTemplate": "download_template",
	}
}

func NewProdukController() *ProdukController {
	return &ProdukController{}
}

// Index renders the product management page
func (ctrl *ProdukController) Index(c *core.Context) error {
	data := core.Map{
		"Title":       "Manajemen Produk",
		"CurrentPage": "produk",
	}

	if err := c.View("inc/header", data); err != nil {
		return err
	}
	if err := c.View("admin/produk/index", data); err != nil {
		return err
	}
	return c.View("inc/footer", data)
}

// GetData handles DataTables server-side AJAX requests
func (ctrl *ProdukController) GetData(c *core.Context) error {
	type ProdukWithKategori struct {
		models.Produk
		NamaKategori string `db:"nama_kategori"`
	}

	var list []ProdukWithKategori
	_ = database.Table("produk").
		Select("produk.*, COALESCE(produk_kategori.nama_kategori, '-') as nama_kategori").
		LeftJoin("produk_kategori", "produk.id_produk_kategori", "=", "produk_kategori.id").
		Where("produk.is_deleted", false).
		OrderBy("produk.id", "ASC").
		Get(&list)

	total, _ := database.Table("produk").Where("is_deleted", false).Count()

	type RowData []any
	var dataRows []RowData

	start, _ := strconv.Atoi(c.Form("start"))
	for i, p := range list {
		actions := fmt.Sprintf(`
			<button class="btn btn-info btn-sm btn-detail" data-id="%d" data-bs-toggle="modal" data-bs-target="#detail-modal">
				<i class="fas fa-eye"></i> Detail
			</button>
			<button class="btn btn-warning btn-sm btn-edit" data-id="%d" data-bs-toggle="modal" data-bs-target="#form-modal">
				<i class="fas fa-edit"></i> Edit
			</button>
			<button class="btn btn-danger btn-sm btn-hapus" data-id="%d">
				<i class="fas fa-trash"></i> Hapus
			</button>
		`, p.ID, p.ID, p.ID)

		status := `<span class="badge bg-danger">Nonaktif</span>`
		if p.IsActive {
			status = `<span class="badge bg-success">Aktif</span>`
		}

		check := fmt.Sprintf(`<input type="checkbox" class="data-check" value="%d">`, p.ID)
		hargaStr := fmt.Sprintf("Rp %s", formatRupiah(p.Harga))

		dataRows = append(dataRows, RowData{
			check,
			start + i + 1,
			p.CodeProduk,
			p.NamaProduk,
			p.NamaKategori,
			hargaStr,
			p.Stok,
			status,
			actions,
		})
	}

	draw, _ := strconv.Atoi(c.Form("draw"))
	return c.JSON(http.StatusOK, core.Map{
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            dataRows,
	})
}

// Tambah loads the create product form modal
func (ctrl *ProdukController) Tambah(c *core.Context) error {
	var kategori []models.ProdukKategori
	_ = database.Table("produk_kategori").Where("is_deleted", false).Get(&kategori)

	data := core.Map{
		"Title":    "Tambah Produk",
		"Action":   "/admin/produk/simpan",
		"Produk":   models.Produk{IsActive: true},
		"Kategori": kategori,
	}

	return c.View("admin/produk/form", data)
}

// Edit loads the edit product form modal
func (ctrl *ProdukController) Edit(c *core.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return c.String(http.StatusBadRequest, "ID tidak valid")
	}

	var produk models.Produk
	if err := database.Table("produk").Where("id", id).First(&produk); err != nil {
		return c.String(http.StatusNotFound, "Produk tidak ditemukan")
	}

	var kategori []models.ProdukKategori
	_ = database.Table("produk_kategori").Where("is_deleted", false).Get(&kategori)

	data := core.Map{
		"Title":    "Edit Produk",
		"Action":   fmt.Sprintf("/admin/produk/update/%d", id),
		"Produk":   produk,
		"Kategori": kategori,
	}

	return c.View("admin/produk/form", data)
}

// Detail loads the product detail modal
func (ctrl *ProdukController) Detail(c *core.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return c.String(http.StatusBadRequest, "ID tidak valid")
	}

	type ProdukDetail struct {
		models.Produk
		NamaKategori string `db:"nama_kategori"`
	}

	var produk ProdukDetail
	if err := database.Table("produk").
		Select("produk.*, COALESCE(produk_kategori.nama_kategori, '-') as nama_kategori").
		LeftJoin("produk_kategori", "produk.id_produk_kategori", "=", "produk_kategori.id").
		Where("produk.id", id).
		First(&produk); err != nil {
		return c.String(http.StatusNotFound, "Produk tidak ditemukan")
	}

	data := core.Map{
		"Produk": produk,
	}

	return c.View("admin/produk/detail", data)
}

// Simpan creates a new product
func (ctrl *ProdukController) Simpan(c *core.Context) error {
	codeProduk := strings.TrimSpace(c.Form("code_produk"))
	namaProduk := strings.TrimSpace(c.Form("nama_produk"))
	idKategori, _ := strconv.Atoi(c.Form("id_produk_kategori"))
	harga, _ := strconv.ParseFloat(c.Form("harga"), 64)
	stok, _ := strconv.Atoi(c.Form("stok"))
	deskripsi := c.Form("deskripsi")
	isActive := c.Form("is_active") == "1"

	if codeProduk == "" || namaProduk == "" || idKategori == 0 {
		return c.JSON(http.StatusBadRequest, core.Map{
			"status":  "error",
			"message": "Validasi gagal: Kode, Nama, dan Kategori wajib diisi",
		})
	}

	err := database.Table("produk").Insert(map[string]any{
		"code_produk":        codeProduk,
		"nama_produk":        namaProduk,
		"id_produk_kategori": idKategori,
		"harga":              harga,
		"stok":               stok,
		"deskripsi":          deskripsi,
		"is_active":          isActive,
		"is_deleted":         false,
		"created_at":         time.Now(),
		"updated_at":         time.Now(),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.Map{
			"status":  "error",
			"message": "Gagal menambahkan produk: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, core.Map{
		"status":  "success",
		"message": "Produk berhasil ditambahkan!",
	})
}

// Update saves changes to an existing product
func (ctrl *ProdukController) Update(c *core.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return c.JSON(http.StatusBadRequest, core.Map{
			"status":  "error",
			"message": "ID tidak valid",
		})
	}

	codeProduk := strings.TrimSpace(c.Form("code_produk"))
	namaProduk := strings.TrimSpace(c.Form("nama_produk"))
	idKategori, _ := strconv.Atoi(c.Form("id_produk_kategori"))
	harga, _ := strconv.ParseFloat(c.Form("harga"), 64)
	stok, _ := strconv.Atoi(c.Form("stok"))
	deskripsi := c.Form("deskripsi")
	isActive := c.Form("is_active") == "1"

	err := database.Table("produk").Where("id", id).Update(map[string]any{
		"code_produk":        codeProduk,
		"nama_produk":        namaProduk,
		"id_produk_kategori": idKategori,
		"harga":              harga,
		"stok":               stok,
		"deskripsi":          deskripsi,
		"is_active":          isActive,
		"updated_at":         time.Now(),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.Map{
			"status":  "error",
			"message": "Gagal memperbarui produk: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, core.Map{
		"status":  "success",
		"message": "Produk berhasil diperbarui!",
	})
}

// Hapus marks a product as deleted (soft delete)
func (ctrl *ProdukController) Hapus(c *core.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return c.JSON(http.StatusBadRequest, core.Map{
			"status":  "error",
			"message": "ID tidak valid",
		})
	}

	err := database.Table("produk").Where("id", id).Update(map[string]any{
		"is_deleted": true,
		"deleted_at": time.Now(),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, core.Map{
			"status":  "error",
			"message": "Gagal menghapus produk: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, core.Map{
		"status":  "success",
		"message": "Produk berhasil dihapus!",
	})
}

// BulkDelete deletes multiple products
func (ctrl *ProdukController) BulkDelete(c *core.Context) error {
	_ = c.Request.ParseForm()
	ids := c.Request.Form["ids[]"]
	if len(ids) == 0 {
		ids = c.Request.Form["ids"]
	}

	if len(ids) == 0 {
		return c.JSON(http.StatusBadRequest, core.Map{
			"status":  "error",
			"message": "Tidak ada produk yang dipilih!",
		})
	}

	for _, idStr := range ids {
		if id, err := strconv.Atoi(idStr); err == nil {
			_ = database.Table("produk").Where("id", id).Update(map[string]any{
				"is_deleted": true,
				"deleted_at": time.Now(),
			})
		}
	}

	return c.JSON(http.StatusOK, core.Map{
		"status":  "success",
		"message": "Produk yang dipilih berhasil dihapus!",
	})
}

// DownloadTemplate provides the CSV template for bulk import
func (ctrl *ProdukController) DownloadTemplate(c *core.Context) error {
	csv := "code_produk,nama_produk,code_kategori,harga,stok,deskripsi\n"
	csv += "AG-001,Ayam Goreng,makanan,13000,100,Contoh deskripsi\n"
	csv += "JZ-001,Jus Jeruk,minuman,10000,50,\n"

	c.Response.Header().Set("Content-Type", "text/csv")
	c.Response.Header().Set("Content-Disposition", "attachment; filename=\"template_import_produk.csv\"")
	return c.String(http.StatusOK, csv)
}

func formatRupiah(val float64) string {
	intVal := int64(val)
	str := strconv.FormatInt(intVal, 10)
	n := len(str)
	if n <= 3 {
		return str
	}
	var res []string
	remainder := n % 3
	if remainder > 0 {
		res = append(res, str[:remainder])
	}
	for i := remainder; i < n; i += 3 {
		res = append(res, str[i:i+3])
	}
	return strings.Join(res, ".")
}
