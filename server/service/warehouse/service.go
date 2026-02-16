package warehouse

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	ListWarehouses() ([]*types.Warehouse, error)
	CreateWarehouse(payload types.WarehousePayload) ([]*types.Warehouse, error)
	UpdateWarehouse(id int, payload types.WarehousePayload) ([]*types.Warehouse, error)
	DeleteWarehouse(id int) ([]*types.Warehouse, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/warehouse", func(rt chi.Router) {
		rt.Get("/", service.HandleGet)
		rt.Post("/", service.HandleCreate)
		rt.Put("/{id}", service.HandleUpdate)
		rt.Delete("/{id}", service.HandleDelete)
	})
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	items, err := s.store.ListWarehouses()
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}

func (s *Service) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	var payload types.WarehousePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.CreateWarehouse(payload)
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
	var payload types.WarehousePayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.UpdateWarehouse(id, payload)
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
	items, err := s.store.DeleteWarehouse(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}
