package middleware

import "net/http"

var payload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (m *Middleware) AuthToken(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.Models.Tokens.Authenticate(r.Header.Get("Authorization"))
		if err != nil {
			payload.Error = true
			payload.Message = "invalid authentication credentials"

			m.App.WriteJson(w, http.StatusUnauthorized, payload)
		}
	})
}
