package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/libraries/database"
	"pos-cafe-go/app/models"
)

type KategoriController struct {
	core.Controller
}

func (ctrl *KategoriController) Routes() map[string]string {
	return map[string]string{
		"GetData": "get_data_kategori",
	}
}

func NewKategoriController() *KategoriController {
	return &KategoriController{}
}

func (ctrl *KategoriController) Index(c *core.Context) error {
	data := core.Map{
		"Title":       "Manajemen Kategori Produk",
		"CurrentPage": "kategori",
	}

	if err := c.View("inc/header", data); err != nil {
		return err
	}
	if err := c.View("admin/kategori/index", data); err != nil {
		return err
	}
	return c.View("inc/footer", data)
}

func (ctrl *KategoriController) GetData(c *core.Context) error {
	var list []models.ProdukKategori
	_ = database.Table("produk_kategori").Where("is_deleted", false).Get(&list)

	total, _ := database.Table("produk_kategori").Where("is_deleted", false).Count()

	type RowData []any
	var dataRows []RowData

	for i, kat := range list {
		actions := fmt.Sprintf(`
			<button class="btn btn-warning btn-sm btn-edit" data-id="%d" data-bs-toggle="modal" data-bs-target="#form-modal">
				<i class="fas fa-edit"></i> Edit
			</button>
			<button class="btn btn-danger btn-sm btn-hapus" data-id="%d">
				<i class="fas fa-trash"></i> Hapus
			</button>
		`, kat.ID, kat.ID)

		check := fmt.Sprintf(`<input type="checkbox" class="data-check" value="%d">`, kat.ID)

		dataRows = append(dataRows, RowData{
			check,
			i + 1,
			kat.CodeKategori,
			kat.NamaKategori,
			kat.CreatedAt.Format("02 Jan 2006 15:04"),
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
