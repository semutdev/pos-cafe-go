package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/helpers"
	"github.com/semutdev/goigniter/system/libraries/database"
	"github.com/semutdev/goigniter/system/libraries/session"
	"pos-cafe-go/app/controllers"
)

func setupTestApp(t *testing.T) *core.Application {
	app := core.New()
	helpers.Init("http://localhost:8080")

	session.Init(session.Config{
		Secret:     "test-secret-key-32-bytes-minimum-length!",
		CookieName: "test_session",
		MaxAge:     86400,
	})

	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	database.SetDefault(db)
	setupDatabase(db)

	funcs := helpers.AllTemplateFuncs()
	funcs["currency"] = func(val any) string { return "" }
	funcs["format_date"] = func(val any) string { return "" }

	if err := app.LoadTemplatesWithFuncs("./app/views", false, funcs); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	core.ResetRegistry()
	app.Register(&controllers.AuthController{})
	app.Register(&controllers.DashboardController{}, "admin")
	app.Register(&controllers.KategoriController{}, "admin")
	app.Register(&controllers.ProdukController{}, "admin")
	app.Register(&controllers.PosController{})
	app.AutoRoute()

	app.GET("/", func(c *core.Context) error {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	})

	return app
}

func TestAppRoutes(t *testing.T) {
	app := setupTestApp(t)

	// Test 1: Root redirect
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("Expected 302 Found, got %d", rec.Code)
	}

	// Test 2: Dashboard via AutoRoute
	req = httptest.NewRequest("GET", "/admin/dashboard", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/dashboard, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Dashboard") {
		t.Errorf("Expected body to contain 'Dashboard'")
	}
	// Verify base_url and site_url helper output (NO double slashes //assets or //admin)
	if strings.Contains(body, `href="//assets/`) || strings.Contains(body, `src="//assets/`) {
		t.Errorf("Found invalid protocol-relative URL //assets in response body!")
	}
	if strings.Contains(body, `href="//admin/`) {
		t.Errorf("Found invalid protocol-relative link //admin in response body!")
	}
	if !strings.Contains(body, "http://localhost:8080/assets/css/adminlte.css") && !strings.Contains(body, "/assets/css/adminlte.css") {
		t.Errorf("Expected base_url css link in body")
	}
	if !strings.Contains(body, "http://localhost:8080/admin/dashboard") && !strings.Contains(body, "/admin/dashboard") {
		t.Errorf("Expected site_url dashboard link in body")
	}
	// Assert main dashboard content is rendered (not empty!)
	if !strings.Contains(body, "app-main") {
		t.Errorf("Expected dashboard body to contain 'app-main' container")
	}
	if !strings.Contains(body, "Produk Stok Menipis") {
		t.Errorf("Expected dashboard body to contain 'Produk Stok Menipis'")
	}
	if !strings.Contains(body, "summary-cards-container") {
		t.Errorf("Expected dashboard body to contain 'summary-cards-container'")
	}

	// Test 2b: Dashboard Index alias (/admin/dashboard/index)
	req = httptest.NewRequest("GET", "/admin/dashboard/index", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/dashboard/index, got %d", rec.Code)
	}

	// Test 3: POS page via AutoRoute
	req = httptest.NewRequest("GET", "/pos", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /pos, got %d", rec.Code)
	}
	posBody := rec.Body.String()
	if !strings.Contains(posBody, "product-grid-container") {
		t.Errorf("Expected /pos body to contain 'product-grid-container'")
	}
	if !strings.Contains(posBody, "cart-content") {
		t.Errorf("Expected /pos body to contain 'cart-content'")
	}

	// Test 3b: POS Index alias (/pos/index)
	req = httptest.NewRequest("GET", "/pos/index", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /pos/index, got %d", rec.Code)
	}

	// Test 4: Kategori page via AutoRoute
	req = httptest.NewRequest("GET", "/admin/kategori", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/kategori, got %d", rec.Code)
	}

	// Test 5: Kategori datatable API via AutoRoute custom Route ("GetData": "get_data_kategori")
	req = httptest.NewRequest("POST", "/admin/kategori/get_data_kategori", strings.NewReader("draw=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for kategori data, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "recordsTotal") {
		t.Errorf("Expected JSON response to contain recordsTotal")
	}

	// Test 6: Produk page via AutoRoute (/admin/produk)
	req = httptest.NewRequest("GET", "/admin/produk", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/produk, got %d", rec.Code)
	}
	produkBody := rec.Body.String()
	if !strings.Contains(produkBody, "Manajemen Produk") {
		t.Errorf("Expected /admin/produk body to contain 'Manajemen Produk'")
	}
	if !strings.Contains(produkBody, "produk-table") {
		t.Errorf("Expected /admin/produk body to contain 'produk-table'")
	}

	// Test 7: Produk datatable API via AutoRoute (/admin/produk/get_data_produk)
	req = httptest.NewRequest("POST", "/admin/produk/get_data_produk", strings.NewReader("draw=1&start=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for produk data, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "recordsTotal") {
		t.Errorf("Expected JSON response to contain recordsTotal")
	}

	// Test 8: Produk form tambah via AutoRoute (/admin/produk/tambah)
	req = httptest.NewRequest("GET", "/admin/produk/tambah", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/produk/tambah, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "produk-form") {
		t.Errorf("Expected tambah form to contain 'produk-form'")
	}

	// Test 9: Produk detail modal via AutoRoute (/admin/produk/detail/1)
	req = httptest.NewRequest("GET", "/admin/produk/detail/1", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /admin/produk/detail/1, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Kode Produk") {
		t.Errorf("Expected detail modal to contain 'Kode Produk'")
	}
}

func TestAppWithMySQL(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/pos-cafe?charset=utf8mb4&parseTime=true"
	db, err := database.Open("mysql", dsn)
	if err != nil {
		t.Skipf("Skipping MySQL test, could not connect: %v", err)
	}
	defer db.Close()
	database.SetDefault(db)

	app := core.New()
	helpers.Init("http://localhost:8080")

	session.Init(session.Config{
		Secret:     "test-secret-key-32-bytes-minimum-length!",
		CookieName: "test_session",
		MaxAge:     86400,
	})

	funcs := helpers.AllTemplateFuncs()
	funcs["currency"] = func(val any) string { return fmt.Sprintf("Rp %v", val) }
	funcs["format_date"] = func(val any) string { return fmt.Sprintf("%v", val) }

	if err := app.LoadTemplatesWithFuncs("./app/views", false, funcs); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	core.ResetRegistry()
	app.Register(&controllers.DashboardController{}, "admin")
	app.Register(&controllers.KategoriController{}, "admin")
	app.Register(&controllers.ProdukController{}, "admin")
	app.Register(&controllers.PosController{})
	app.AutoRoute()

	// 1. Dashboard with MySQL
	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 for dashboard with MySQL, got %d", rec.Code)
	}
	t.Logf("MySQL Dashboard rendered %d bytes", rec.Body.Len())

	// 2. Produk list with MySQL
	req = httptest.NewRequest("GET", "/admin/produk", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 for /admin/produk with MySQL, got %d", rec.Code)
	}
	t.Logf("MySQL Produk page rendered %d bytes", rec.Body.Len())

	// 3. Produk datatable API with MySQL
	req = httptest.NewRequest("POST", "/admin/produk/get_data_produk", strings.NewReader("draw=1&start=0"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 for /admin/produk/get_data_produk with MySQL, got %d", rec.Code)
	}
	t.Logf("MySQL Produk DataTable response: %s", rec.Body.String())

	// 4. POS page with MySQL
	req = httptest.NewRequest("GET", "/pos", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 for /pos with MySQL, got %d", rec.Code)
	}
	t.Logf("MySQL POS page rendered %d bytes", rec.Body.Len())
}

func TestAuthRoutes(t *testing.T) {
	app := setupTestApp(t)

	// Test 1: GET /auth/logout should redirect to /auth/login (NOT 404!)
	req := httptest.NewRequest("GET", "/auth/logout", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("Expected 302 Found for /auth/logout, got %d", rec.Code)
	}
	location := rec.Header().Get("Location")
	if location != "/auth/login" {
		t.Errorf("Expected redirect location to be /auth/login, got %s", location)
	}

	// Test 2: GET /auth/login should return 200 OK
	req = httptest.NewRequest("GET", "/auth/login", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for /auth/login, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Masukan email dan sandi!") {
		t.Errorf("Expected login page body to contain 'Masukan email dan sandi!'")
	}

	// Test 3: POST /auth/login with valid credentials redirects to /admin/dashboard
	formBody := strings.NewReader("identity=admin@admin.com&password=password")
	req = httptest.NewRequest("POST", "/auth/login", formBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("Expected 302 Found on successful login, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/admin/dashboard" {
		t.Errorf("Expected redirect to /admin/dashboard, got %s", rec.Header().Get("Location"))
	}

	// Test 4: POST /auth/login with invalid password redirects back to /auth/login
	formBody = strings.NewReader("identity=admin@admin.com&password=wrongpassword")
	req = httptest.NewRequest("POST", "/auth/login", formBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("Expected 302 Found on failed login, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/auth/login" {
		t.Errorf("Expected redirect to /auth/login, got %s", rec.Header().Get("Location"))
	}
}

