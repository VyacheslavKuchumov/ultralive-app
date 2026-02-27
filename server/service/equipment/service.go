package equipment

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	SearchEquipment(query types.ListQuery) ([]*types.Equipment, int, error)
	SearchEquipmentBySetID(setID int, query types.ListQuery) ([]*types.Equipment, int, error)
	ListEquipment() ([]*types.Equipment, error)
	ListEquipmentBySetID(setID int) ([]*types.Equipment, error)
	GetEquipmentByID(id int) (*types.Equipment, error)
	CreateEquipment(payload types.EquipmentPayload) ([]*types.Equipment, error)
	UpdateEquipment(id int, payload types.EquipmentPayload) ([]*types.Equipment, error)
	DeleteEquipment(id int) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/equipment", func(rt chi.Router) {
		rt.Get("/", service.HandleGet)
		rt.Get("/set/{id}", service.HandleGetBySetID)
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
	items, total, err := s.store.SearchEquipment(query)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.NewPaginatedResponse(items, query.Page, query.PerPage, total))
}

func (s *Service) HandleGetBySetID(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	id, ok := crmhttp.MustPathID(w, r, "id")
	if !ok {
		return
	}
	query := crmhttp.ParseListQuery(r)
	items, total, err := s.store.SearchEquipmentBySetID(id, query)
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
	item, err := s.store.GetEquipmentByID(id)
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
	var payload types.EquipmentPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.CreateEquipment(payload)
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
	var payload types.EquipmentPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.UpdateEquipment(id, payload)
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
	if err := s.store.DeleteEquipment(id); err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
