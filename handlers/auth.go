package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"myapp/data"
	"net/http"
	"time"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "auth/login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Info("Error rendering login page: ", err)
		h.App.ErrInternalServer(w, r)
	}
}

func (h *Handlers) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	pass := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	matches, err := user.PasswordMatches(pass)
	if err != nil {
		w.Write([]byte("Error validating password"))
	}

	if !matches {
		w.Write([]byte("Invalid password"))
	}

	if r.Form.Get("remember") == "remember" {
		randomString := h.App.RandomString(12)
		hasher := sha256.New()
		_, err := hasher.Write([]byte(randomString))
		if err != nil {
			h.App.ErrStatus(w, http.StatusBadRequest)
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		rm := data.RememberToken{}
		err = rm.InsertToken(user.ID, sha)
		if err != nil {
			return
		}

		expiresAt := time.Now().Add(2 * 24 * 60 * 60 * time.Second)
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_remember", h.App.AppName),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha),
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			Domain:   h.App.Session.Cookie.Domain,
			MaxAge:   1000,
			Secure:   h.App.Session.Cookie.Secure,
			SameSite: http.SameSiteDefaultMode,
		}

		http.SetCookie(w, &cookie)
		h.App.Session.Put(r.Context(), "remember_token", sha)
	}

	h.App.Session.Put(r.Context(), "UserId", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if h.App.Session.Exists(r.Context(), "remember_token") {
		rt := data.RememberToken{}
		_ = rt.Delete(h.App.Session.GetString(r.Context(), "remember_token"))
	}

	// delete cookie
	cookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", h.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Domain:   h.App.Session.Cookie.Domain,
		Secure:   h.App.Session.Cookie.Secure,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)

	h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "UserId")
	h.App.Session.Remove(r.Context(), "remember_token")
	h.App.Session.Destroy(r.Context())
	h.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (h *Handlers) Forgot(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "auth/forgot", nil, nil)
	if err != nil {
		h.App.ErrorLog.Error("Error rendering forget page: ", err)
		h.App.ErrInternalServer(w, r)
	}
}

func (h *Handlers) PostForgot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.ErrStatus(w, http.StatusBadRequest)
		return
	}
	var u *data.User
	email := r.Form.Get("email")
	u, err = u.GetByEmail(email)
	if err != nil {
		h.App.ErrStatus(w, http.StatusBadRequest)
		return
	}

	w.Write([]byte("Email functinality is not provided yet by the framework"))
}
