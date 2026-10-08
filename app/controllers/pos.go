package controllers

import (
	"strconv"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/libraries/database"
	"github.com/semutdev/goigniter/system/libraries/session"
	"pos-cafe-go/app/models"
)

type PosController struct {
	core.Controller
}

func NewPosController() *PosController {
	return &PosController{}
}

type CartItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

func (ctrl *PosController) Index(c *core.Context) error {
	var listProduk []models.Produk
	_ = database.Table("produk").Where("is_active", true).Where("is_deleted", false).Get(&listProduk)

	var listKategori []models.ProdukKategori
	_ = database.Table("produk_kategori").Where("is_deleted", false).Get(&listKategori)

	sess := session.Get(c)
	cart := sess.Get("cart")
	if cart == nil {
		cart = []CartItem{}
	}

	data := core.Map{
		"Title":        "Point of Sale",
		"ListProduk":   listProduk,
		"ListKategori": listKategori,
		"Cart":         cart,
	}

	if err := c.View("inc/header", data); err != nil {
		return err
	}
	if err := c.View("pos/index", data); err != nil {
		return err
	}
	return c.View("inc/footer", data)
}

func (ctrl *PosController) FilterProduk(c *core.Context) error {
	keyword := c.Form("keyword")
	idKategori := c.Form("id_kategori")

	query := database.Table("produk").Where("is_active", true).Where("is_deleted", false)
	if keyword != "" {
		query = query.Where("nama_produk", "LIKE", "%"+keyword+"%")
	}
	if idKategori != "" && idKategori != "all" {
		if id, err := strconv.Atoi(idKategori); err == nil {
			query = query.Where("id_produk_kategori", id)
		}
	}

	var listProduk []models.Produk
	_ = query.Get(&listProduk)

	data := core.Map{
		"ListProduk": listProduk,
	}

	return c.View("pos/partials/_product_grid", data)
}
