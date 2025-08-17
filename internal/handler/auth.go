package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/quiby-ai/auth-service/internal/models"
	"github.com/quiby-ai/common/pkg/auth"
)

func LoginWithTelegram(cfg *auth.JWTConfig, userRepo *models.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tgID := user.ID

		profBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "profile marshal error", http.StatusInternalServerError)
			return
		}

		err = userRepo.UpsertUser(r.Context(), tgID, profBytes)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		userIdentity := auth.UserIdentity{
			UserID: strconv.FormatInt(user.ID, 10),
		}
		token, err := auth.IssueAccessJWT(userIdentity, cfg)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]any{
			"user_id": userIdentity.UserID,
			"ok":      true,
			"token":   token,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "response encoding error", http.StatusInternalServerError)
			return
		}
	}
}

func Me(userRepo *models.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tgID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			http.Error(w, "bad user id", http.StatusBadRequest)
			return
		}

		profile, err := userRepo.GetUserProfile(r.Context(), tgID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		if profile == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(profile); err != nil {
			http.Error(w, "response write error", http.StatusInternalServerError)
			return
		}
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		http.Error(w, "response write error", http.StatusInternalServerError)
		return
	}
}
