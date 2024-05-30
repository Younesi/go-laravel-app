package middleware

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Middleware) CheckRemember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.App.Session.Exists(r.Context(), "UserId") {
			cookie, err := r.Cookie(fmt.Sprintf("_%s_remember", m.App.AppName))
			if err == nil {
				key := cookie.Value
				var u data.User
				if len(key) > 0 {
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := u.CheckForRememberToken(id, hash)
					if !validHash {
						m.deleteRememberCookie(w, r)
						m.App.Session.Put(r.Context(), "error", "You've been logged out from another device")
					} else {
						user, _ := u.Get(id)
						m.App.Session.Put(r.Context(), "UserId", user.ID)
						m.App.Session.Put(r.Context(), "remember_token", hash)
					}
				} else {
					m.deleteRememberCookie(w, r)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) deleteRememberCookie(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", m.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   m.App.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   m.App.Session.Cookie.Secure,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &newCookie)

	m.App.Session.Remove(r.Context(), "UserId")
	m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
}
