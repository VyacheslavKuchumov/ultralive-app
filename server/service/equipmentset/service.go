package equipmentset

import (
	"VyacheslavKuchumov/test-backend/service/crmhttp"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	SearchEquipmentSets(query types.ListQuery) ([]*types.EquipmentSet, int, error)
	ListEquipmentSets() ([]*types.EquipmentSet, error)
	GetEquipmentSetByID(id int) (*types.EquipmentSet, error)
	CreateEquipmentSet(payload types.EquipmentSetPayload) ([]*types.EquipmentSet, error)
	UpdateEquipmentSet(id int, payload types.EquipmentSetPayload) ([]*types.EquipmentSet, error)
	DeleteEquipmentSet(id int) ([]*types.EquipmentSet, error)
	GetEquipmentSetsWithMaintenance() ([]*types.EquipmentSet, error)
	GetEquipmentSetsWithStorage() ([]*types.EquipmentSetStorageSummary, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func RegisterRoutes(r chi.Router, service *Service) {
	r.Route("/equipment_set", func(rt chi.Router) {
		rt.Get("/", service.HandleGet)
		rt.Get("/search/{id}", service.HandleGetByID)
		rt.Get("/maintenance", service.HandleGetWithMaintenance)
		rt.Get("/storage", service.HandleGetWithStorage)
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
	items, total, err := s.store.SearchEquipmentSets(query)
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
	item, err := s.store.GetEquipmentSetByID(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, item)
}

func (s *Service) HandleGetWithMaintenance(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	items, err := s.store.GetEquipmentSetsWithMaintenance()
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}

func (s *Service) HandleGetWithStorage(w http.ResponseWriter, r *http.Request) {
	if !crmhttp.RequireAuth(w, r) {
		return
	}
	items, err := s.store.GetEquipmentSetsWithStorage()
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
	var payload types.EquipmentSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.CreateEquipmentSet(payload)
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
	var payload types.EquipmentSetPayload
	if !crmhttp.ParseAndValidate(w, r, &payload) {
		return
	}
	items, err := s.store.UpdateEquipmentSet(id, payload)
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
	items, err := s.store.DeleteEquipmentSet(id)
	if err != nil {
		crmhttp.WriteStoreError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}
