package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/helpers"
	"github.com/semutdev/goigniter/system/libraries/database"
	"github.com/semutdev/goigniter/system/libraries/session"
	"github.com/semutdev/goigniter/system/middleware"
	"pos-cafe-go/app/controllers"
)

func main() {
	loadEnv(".env")

	app := core.New()

	// 1. Initialize Helpers
	helpers.Init("http://localhost:8080")

	// 2. Initialize Session
	session.Init(session.Config{
		Secret:     "pos-cafe-secret-key-must-be-32-bytes-long",
		CookieName: "pos_cafe_session",
		MaxAge:     86400 * 7,
	})

	// 3. Initialize Database (MySQL default)
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "mysql"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pos-cafe"
	}

	var dsn string
	if dbDriver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	} else {
		dsn = "pos_cafe.db"
	}

	db, err := database.Open(dbDriver, dsn)
	if err != nil {
		log.Printf("[DB WARNING] Failed to connect to %s (%s): %v, falling back to SQLite", dbDriver, dsn, err)
		dbDriver = "sqlite"
		dsn = "pos_cafe.db"
		db, err = database.Open("sqlite", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		setupDatabase(db)
	} else if dbDriver == "sqlite" {
		setupDatabase(db)
	}
	database.SetDefault(db)

	// 4. Load Templates with Helper Functions
	funcs := helpers.AllTemplateFuncs()
	funcs["currency"] = func(val any) string {
		return fmt.Sprintf("Rp %v", val)
	}
	funcs["format_date"] = func(val any) string {
		return fmt.Sprintf("%v", val)
	}

	if err := app.LoadTemplatesWithFuncs("./app/views", true, funcs); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// 5. Global Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	// 6. Static Asset Serving
	app.Static("/assets/", "./public/assets")
	app.Static("/uploads/", "./public/uploads")
	app.Static("/public/", "./public")

	// 7. AutoRoute Controllers Registration
	app.Register(&controllers.AuthController{})
	app.Register(&controllers.DashboardController{}, "admin")
	app.Register(&controllers.KategoriController{}, "admin")
	app.Register(&controllers.ProdukController{}, "admin")
	app.Register(&controllers.PosController{})
	app.AutoRoute()

	// 8. Root Route
	app.GET("/", func(c *core.Context) error {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Println("==================================================")
	fmt.Println("  POS Cafe (GoIgniter Framework Edition)")
	fmt.Println("  Successfully migrated from CodeIgniter 3")
	fmt.Println("==================================================")
	fmt.Printf("  Server listening on http://localhost%s\n", port)
	fmt.Printf("  Database   : %s (%s)\n", strings.ToUpper(dbDriver), dbName)
	fmt.Println("  - Dashboard: http://localhost:8080/admin/dashboard")
	fmt.Println("  - POS App  : http://localhost:8080/pos")
	fmt.Println("  - Kategori : http://localhost:8080/admin/kategori")
	fmt.Println("  - Produk   : http://localhost:8080/admin/produk")
	fmt.Println("  - Login    : http://localhost:8080/auth/login")
	fmt.Println("  - Logout   : http://localhost:8080/auth/logout")
	fmt.Println("==================================================")

	log.Fatal(app.Run(port))
}

func setupDatabase(db *database.DB) {
	// Create tables if not exist
	schema := `
	CREATE TABLE IF NOT EXISTS produk_kategori (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code_kategori TEXT NOT NULL UNIQUE,
		nama_kategori TEXT NOT NULL,
		deskripsi TEXT,
		is_deleted INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS produk (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		id_produk_kategori INTEGER NOT NULL,
		code_produk TEXT NOT NULL UNIQUE,
		nama_produk TEXT NOT NULL,
		deskripsi TEXT,
		harga REAL NOT NULL DEFAULT 0,
		stok INTEGER NOT NULL DEFAULT 0,
		gambar TEXT,
		is_active INTEGER DEFAULT 1,
		is_deleted INTEGER DEFAULT 0,
		deleted_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS transaksi (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nomor_transaksi TEXT NOT NULL UNIQUE,
		tanggal DATETIME DEFAULT CURRENT_TIMESTAMP,
		total_akhir REAL NOT NULL DEFAULT 0,
		catatan TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		activation_selector TEXT,
		activation_code TEXT,
		forgotten_password_selector TEXT,
		forgotten_password_code TEXT,
		forgotten_password_time INTEGER,
		remember_selector TEXT,
		remember_code TEXT,
		created_on INTEGER,
		last_login INTEGER,
		active INTEGER DEFAULT 1,
		first_name TEXT,
		last_name TEXT,
		company TEXT,
		phone TEXT
	);
	`
	_, _ = db.Exec(schema)

	// Seed users if empty
	userCount, _ := db.Table("users").Count()
	if userCount == 0 {
		_ = db.Table("users").Insert(map[string]any{
			"username":   "administrator",
			"email":      "admin@admin.com",
			"password":   "$2y$10$jAw6X2eZ99TwPympOZDdFOBY5aH1diznVIiOPpl9TolslxdnKBUOO", // password
			"first_name": "Admin",
			"last_name":  "POS Cafe",
			"active":     1,
		})
	}

	// Seed if empty
	count, _ := db.Table("produk_kategori").Count()
	if count == 0 {
		_ = db.Table("produk_kategori").Insert(map[string]any{
			"code_kategori": "KAT-001",
			"nama_kategori": "Coffee",
			"deskripsi":     "Aneka racikan kopi pilihan",
			"is_deleted":    0,
		})
		_ = db.Table("produk_kategori").Insert(map[string]any{
			"code_kategori": "KAT-002",
			"nama_kategori": "Non-Coffee",
			"deskripsi":     "Minuman segar non-kopi",
			"is_deleted":    0,
		})
		_ = db.Table("produk_kategori").Insert(map[string]any{
			"code_kategori": "KAT-003",
			"nama_kategori": "Snacks & Bakery",
			"deskripsi":     "Kudapan pendamping santai",
			"is_deleted":    0,
		})
	}

	prodCount, _ := db.Table("produk").Count()
	if prodCount == 0 {
		_ = db.Table("produk").Insert(map[string]any{
			"id_produk_kategori": 1,
			"code_produk":        "PRD-001",
			"nama_produk":        "Kopi Susu Gula Aren",
			"deskripsi":          "Espresso signature dengan gula aren asli",
			"harga":              22000,
			"stok":               35,
			"is_active":          1,
			"is_deleted":         0,
		})
		_ = db.Table("produk").Insert(map[string]any{
			"id_produk_kategori": 1,
			"code_produk":        "PRD-002",
			"nama_produk":        "Americano Iced",
			"deskripsi":          "Double shot espresso dengan es segar",
			"harga":              18000,
			"stok":               50,
			"is_active":          1,
			"is_deleted":         0,
		})
		_ = db.Table("produk").Insert(map[string]any{
			"id_produk_kategori": 2,
			"code_produk":        "PRD-003",
			"nama_produk":        "Matcha Latte",
			"deskripsi":          "Pure Uji Matcha dengan susu segar",
			"harga":              25000,
			"stok":               8, // stok menipis (<10) untuk tes widget dashboard!
			"is_active":          1,
			"is_deleted":         0,
		})
		_ = db.Table("produk").Insert(map[string]any{
			"id_produk_kategori": 3,
			"code_produk":        "PRD-004",
			"nama_produk":        "Butter Croissant",
			"deskripsi":          "Flaky and buttery french pastry",
			"harga":              18000,
			"stok":               5, // stok menipis (<10)
			"is_active":          1,
			"is_deleted":         0,
		})
	}
}

func loadEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if os.Getenv(k) == "" {
				_ = os.Setenv(k, v)
			}
		}
	}
}
