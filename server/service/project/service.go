package project

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	SearchProjects(archived bool, query types.ListQuery) ([]*types.Project, int, error)
	ListProjects(archived bool) ([]*types.Project, error)
	GetProjectByID(id int) (*types.Project, error)
	CreateProject(payload types.ProjectPayload) ([]*types.Project, error)
	UpdateProject(id int, payload types.ProjectPayload) ([]*types.Project, error)
	DeleteProject(id int) ([]*types.Project, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/projects", func(rt chi.Router) {
		rt.Get("/", service.HandleGet)
		rt.Get("/archived", service.HandleGetArchived)
		rt.Get("/search/{id}", service.HandleGetByID)
		rt.Post("/", service.HandleCreate)
		rt.Put("/{id}", service.HandleUpdate)
		rt.Delete("/{id}", service.HandleDelete)
	})
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	query := crmhttp.ParseListQuery(r)
	items, total, err := s.store.SearchProjects(false, query)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.NewPaginatedResponse(items, query.Page, query.PerPage, total))
}

func (s *Service) HandleGetArchived(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	query := crmhttp.ParseListQuery(r)
	items, total, err := s.store.SearchProjects(true, query)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.NewPaginatedResponse(items, query.Page, query.PerPage, total))
}

func (s *Service) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	id, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	item, err := s.store.GetProjectByID(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, item)
}

func (s *Service) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.ProjectPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.CreateProject(payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, items)
}

func (s *Service) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	id, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	var payload types.ProjectPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.UpdateProject(id, payload)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}

func (s *Service) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	id, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	items, err := s.store.DeleteProject(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}
