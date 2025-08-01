package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clientpulse-org/common/pkg/auth"
)

func LoginWithTelegram(jwtCfg *auth.JWTConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := auth.IssueJWTFromTelegramUser(user, jwtCfg)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		response := map[string]any{
			"token": token,
			"user":  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
