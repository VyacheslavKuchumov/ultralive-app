package equipmentindraft

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	GetEquipmentInDraft(draftID int) (*types.EquipmentInDraftResponse, error)
	AddEquipmentToDraft(payload types.EquipmentInDraftPayload) (*types.EquipmentInDraftResponse, error)
	RemoveEquipmentFromDraft(payload types.DraftEquipmentDeletePayload) (*types.EquipmentInDraftResponse, error)
	AddSetToDraft(payload types.DraftSetPayload) (*types.EquipmentInDraftResponse, error)
	RemoveSetFromDraft(payload types.DraftSetDeletePayload) (*types.EquipmentInDraftResponse, error)
	GetAvailableDraftEquipmentInSet(payload types.DraftSetPayload) ([]*types.Equipment, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/equipment_in_draft", func(rt chi.Router) {
		rt.Get("/{id}", service.HandleGet)
		rt.Post("/add", service.HandleAddEquipment)
		rt.Put("/del", service.HandleRemoveEquipment)
		rt.Post("/add_set", service.HandleAddSet)
		rt.Put("/del_set", service.HandleRemoveSet)
		rt.Post("/equipment_in_set", service.HandleGetAvailableInSet)
	})
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	draftID, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	response, err := s.store.GetEquipmentInDraft(draftID)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.EquipmentInDraftPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.AddEquipmentToDraft(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleRemoveEquipment(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.DraftEquipmentDeletePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.RemoveEquipmentFromDraft(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleAddSet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.DraftSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.AddSetToDraft(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleRemoveSet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.DraftSetDeletePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.RemoveSetFromDraft(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleGetAvailableInSet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.DraftSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.GetAvailableDraftEquipmentInSet(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
