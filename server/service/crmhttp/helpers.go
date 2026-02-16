package crmhttp

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/service/tracker"
	"VyacheslavKuchumov/test-backend/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RequireAuth(w http.ResponseWriter, r *http.Request) bool {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return false
	}
	return true
}

func ParseAndValidate(w http.ResponseWriter, r *http.Request, payload any) bool {
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return false
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrs, ok := err.(validator.ValidationErrors)
		if ok {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrs))
			return false
		}
		utils.WriteError(w, http.StatusBadRequest, err)
		return false
	}

	return true
}

func MustPathID(w http.ResponseWriter, r *http.Request, key string) (int, bool) {
	idRaw := chi.URLParam(r, key)
	id, err := strconv.Atoi(idRaw)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return 0, false
	}
	return id, true
}

func WriteStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, tracker.ErrNotFound):
		utils.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, tracker.ErrInvalidReference):
		utils.WriteError(w, http.StatusBadRequest, err)
	default:
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
