package settype

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	SearchSetTypes(query types.ListQuery) ([]*types.SetType, int, error)
	ListSetTypes() ([]*types.SetType, error)
	GetSetTypeByID(id int) (*types.SetType, error)
	CreateSetType(payload types.SetTypePayload) ([]*types.SetType, error)
	UpdateSetType(id int, payload types.SetTypePayload) ([]*types.SetType, error)
	DeleteSetType(id int) ([]*types.SetType, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/set_types", func(rt chi.Router) {
		rt.Get("/", service.HandleGet)
		rt.Get("/{id}", service.HandleGetByID)
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
	items, total, err := s.store.SearchSetTypes(query)
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
	item, err := s.store.GetSetTypeByID(id)
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
	var payload types.SetTypePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.CreateSetType(payload)
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
	var payload types.SetTypePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.UpdateSetType(id, payload)
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
	items, err := s.store.DeleteSetType(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}
