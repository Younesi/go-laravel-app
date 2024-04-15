package handlers

import "net/http"

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "auth/login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
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

	h.App.Session.Put(r.Context(), "UserId", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "UserId")

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
