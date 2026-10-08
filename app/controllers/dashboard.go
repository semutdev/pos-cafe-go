package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/libraries/database"
	"pos-cafe-go/app/models"
)

type DashboardController struct {
	core.Controller
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (ctrl *DashboardController) Index(c *core.Context) error {
	today := time.Now().Format("2006-01-02")

	// Total Produk
	totalProduk, _ := database.Table("produk").Where("is_deleted", false).Count()

	// Transaction summary today
	var summary struct {
		JumlahTransaksi int     `db:"jumlah_transaksi"`
		TotalPendapatan float64 `db:"total_pendapatan"`
	}
	_ = database.Table("transaksi").
		Select("COUNT(id) as jumlah_transaksi, COALESCE(SUM(total_akhir), 0) as total_pendapatan").
		Where("tanggal", ">=", today+" 00:00:00").
		Where("tanggal", "<=", today+" 23:59:59").
		First(&summary)

	// Low stock products (< 10)
	var lowStock []models.Produk
	_ = database.Table("produk").
		Where("stok", "<=", 10).
		Where("is_deleted", false).
		OrderBy("stok", "ASC").
		Get(&lowStock)

	data := core.Map{
		"Title":           "Dashboard",
		"TotalProduk":     totalProduk,
		"JumlahTransaksi": summary.JumlahTransaksi,
		"TotalPendapatan": summary.TotalPendapatan,
		"LowStockProduk":  lowStock,
		"ActiveFilter":    "today",
	}

	if err := c.View("inc/header", data); err != nil {
		return err
	}
	if err := c.View("admin/dashboard/index", data); err != nil {
		return err
	}
	return c.View("inc/footer", data)
}

func (ctrl *DashboardController) FilterSummary(c *core.Context) error {
	rangeParam := c.Query("range")
	now := time.Now()
	startDate := now.Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	switch rangeParam {
	case "this_week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		end := start.AddDate(0, 0, 6)
		startDate = start.Format("2006-01-02")
		endDate = end.Format("2006-01-02")
	case "this_month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = lastDay
	}

	var summary struct {
		JumlahTransaksi int     `db:"jumlah_transaksi"`
		TotalPendapatan float64 `db:"total_pendapatan"`
	}
	_ = database.Table("transaksi").
		Select("COUNT(id) as jumlah_transaksi, COALESCE(SUM(total_akhir), 0) as total_pendapatan").
		Where("tanggal", ">=", startDate+" 00:00:00").
		Where("tanggal", "<=", endDate+" 23:59:59").
		First(&summary)

	data := core.Map{
		"JumlahTransaksi": summary.JumlahTransaksi,
		"TotalPendapatan": fmt.Sprintf("%.0f", summary.TotalPendapatan),
	}

	return c.ViewWithCode(http.StatusOK, "admin/dashboard/partials/_summary_cards", data)
}
