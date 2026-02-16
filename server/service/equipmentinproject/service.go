package equipmentinproject

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	GetEquipmentInProject(projectID int) (*types.EquipmentInProjectResponse, error)
	AddEquipmentToProject(payload types.EquipmentInProjectPayload) (*types.EquipmentInProjectResponse, error)
	RemoveEquipmentFromProject(payload types.ProjectEquipmentDeletePayload) (*types.EquipmentInProjectResponse, error)
	AddSetToProject(payload types.ProjectSetPayload) (*types.EquipmentInProjectResponse, error)
	RemoveSetFromProject(payload types.ProjectSetDeletePayload) (*types.EquipmentInProjectResponse, error)
	GetAvailableProjectEquipmentInSet(payload types.ProjectSetPayload) ([]*types.Equipment, error)
	GetConflictingEquipment(projectID int) ([]*types.EquipmentConflict, error)
	AddDraftToProject(payload types.AddDraftToProjectPayload) (*types.EquipmentInProjectResponse, error)
	ResetEquipmentInProject(projectID int) error
	GetConflictingProjects() ([]*types.ConflictingProject, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/equipment_in_project", func(rt chi.Router) {
		rt.Get("/{id}", service.HandleGet)
		rt.Post("/add", service.HandleAddEquipment)
		rt.Put("/del", service.HandleRemoveEquipment)
		rt.Post("/add_set", service.HandleAddSet)
		rt.Put("/del_set", service.HandleRemoveSet)
		rt.Post("/equipment_in_set", service.HandleGetAvailableInSet)
		rt.Post("/conflicting", service.HandleGetConflictingEquipment)
		rt.Delete("/reset/{id}", service.HandleReset)
		rt.Post("/add_draft", service.HandleAddDraft)
		rt.Post("/conflicting_projects", service.HandleGetConflictingProjects)
	})
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	projectID, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	response, err := s.store.GetEquipmentInProject(projectID)
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
	var payload types.EquipmentInProjectPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.AddEquipmentToProject(payload)
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
	var payload types.ProjectEquipmentDeletePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.RemoveEquipmentFromProject(payload)
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
	var payload types.ProjectSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.AddSetToProject(payload)
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
	var payload types.ProjectSetDeletePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.RemoveSetFromProject(payload)
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
	var payload types.ProjectSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.GetAvailableProjectEquipmentInSet(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleGetConflictingEquipment(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload struct {
		ProjectID int `json:"project_id" validate:"required,min=1"`
	}
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.GetConflictingEquipment(payload.ProjectID)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleReset(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	projectID, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	if err := s.store.ResetEquipmentInProject(projectID); err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "equipment reset"})
}

func (s *Service) HandleAddDraft(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.AddDraftToProjectPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	response, err := s.store.AddDraftToProject(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Service) HandleGetConflictingProjects(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	response, err := s.store.GetConflictingProjects()
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
