package controllers

import (
	"net/http"
	"strings"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/libraries/database"
	"github.com/semutdev/goigniter/system/libraries/session"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	core.Controller
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Index redirects to login or dashboard
func (ctrl *AuthController) Index(c *core.Context) error {
	sess := session.Get(c)
	if sess.Get("user_id") != nil {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}
	return c.Redirect(http.StatusFound, "/auth/login")
}

// Login handles GET (show form) and POST (process login)
func (ctrl *AuthController) Login(c *core.Context) error {
	if c.Method() == http.MethodPost {
		identity := strings.TrimSpace(c.FormValue("identity"))
		password := c.FormValue("password")

		if identity == "" || password == "" {
			session.SetFlash(c, "message", "Email/username dan password wajib diisi.")
			return c.Redirect(http.StatusFound, "/auth/login")
		}

		// Check user in database
		var user struct {
			ID        int     `db:"id"`
			Username  string  `db:"username"`
			Password  string  `db:"password"`
			Email     string  `db:"email"`
			FirstName *string `db:"first_name"`
		}
		err := database.Table("users").
			Select("id", "username", "password", "email", "first_name").
			Where("email", identity).
			OrWhere("username", identity).
			First(&user)

		if err != nil || user.ID == 0 {
			session.SetFlash(c, "message", "Pengguna tidak ditemukan atau kredensial salah.")
			return c.Redirect(http.StatusFound, "/auth/login")
		}

		// Verify password using bcrypt (normalizing PHP $2y$ prefix to standard $2a$)
		hash := user.Password
		if strings.HasPrefix(hash, "$2y$") {
			hash = "$2a$" + hash[4:]
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			// Fallback check if plain text matches
			if user.Password != password {
				session.SetFlash(c, "message", "Kata sandi salah.")
				return c.Redirect(http.StatusFound, "/auth/login")
			}
		}

		// Authentication successful
		sess := session.Get(c)
		sess.Set("user_id", user.ID)
		sess.Set("email", user.Email)
		sess.Set("username", user.Username)
		if user.FirstName != nil {
			sess.Set("first_name", *user.FirstName)
		}
		if err := sess.Save(c); err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}

	// GET request: show login page
	sess := session.Get(c)
	if sess.Get("user_id") != nil {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}

	message := session.GetFlash(c, "message")
	data := core.Map{
		"Title":   "Login - POS Cafe",
		"Message": message,
	}
	return c.View("auth/login", data)
}

// Logout logs out the current user and redirects to /auth/login
func (ctrl *AuthController) Logout(c *core.Context) error {
	sess := session.Get(c)
	sess.Destroy(c)
	return c.Redirect(http.StatusFound, "/auth/login")
}
