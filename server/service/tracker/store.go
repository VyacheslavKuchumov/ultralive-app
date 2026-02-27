package tracker

import (
	"VyacheslavKuchumov/test-backend/types"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidReference = errors.New("invalid reference")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListSetTypes() ([]*types.SetType, error) {
	rows, err := s.db.Query(`SELECT set_type_id, set_type_name FROM set_types ORDER BY set_type_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.SetType, 0)
	for rows.Next() {
		item := new(types.SetType)
		if err := rows.Scan(&item.SetTypeID, &item.SetTypeName); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) SearchSetTypes(query types.ListQuery) ([]*types.SetType, int, error) {
	items, err := s.ListSetTypes()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.SetType, 0, len(items))
	for _, item := range items {
		if matchesSearch(search, strconv.Itoa(item.SetTypeID), item.SetTypeName) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetSetTypeByID(id int) (*types.SetType, error) {
	row := s.db.QueryRow(`SELECT set_type_id, set_type_name FROM set_types WHERE set_type_id = $1`, id)
	item := new(types.SetType)
	if err := row.Scan(&item.SetTypeID, &item.SetTypeName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (s *Store) CreateSetType(payload types.SetTypePayload) ([]*types.SetType, error) {
	_, err := s.db.Exec(`INSERT INTO set_types (set_type_name) VALUES ($1)`, payload.SetTypeName)
	if err != nil {
		return nil, err
	}
	return s.ListSetTypes()
}

func (s *Store) UpdateSetType(id int, payload types.SetTypePayload) ([]*types.SetType, error) {
	result, err := s.db.Exec(`UPDATE set_types SET set_type_name = $1 WHERE set_type_id = $2`, payload.SetTypeName, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListSetTypes()
}

func (s *Store) DeleteSetType(id int) ([]*types.SetType, error) {
	result, err := s.db.Exec(`DELETE FROM set_types WHERE set_type_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListSetTypes()
}

func (s *Store) ListProjectTypes() ([]*types.ProjectType, error) {
	rows, err := s.db.Query(`SELECT project_type_id, project_type_name, COALESCE(neaktor_id, '') FROM project_types ORDER BY project_type_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.ProjectType, 0)
	for rows.Next() {
		item := new(types.ProjectType)
		if err := rows.Scan(&item.ProjectTypeID, &item.ProjectTypeName, &item.NeaktorID); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) SearchProjectTypes(query types.ListQuery) ([]*types.ProjectType, int, error) {
	items, err := s.ListProjectTypes()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.ProjectType, 0, len(items))
	for _, item := range items {
		if matchesSearch(search, strconv.Itoa(item.ProjectTypeID), item.ProjectTypeName, item.NeaktorID) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetProjectTypeByID(id int) (*types.ProjectType, error) {
	row := s.db.QueryRow(`SELECT project_type_id, project_type_name, COALESCE(neaktor_id, '') FROM project_types WHERE project_type_id = $1`, id)
	item := new(types.ProjectType)
	if err := row.Scan(&item.ProjectTypeID, &item.ProjectTypeName, &item.NeaktorID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (s *Store) CreateProjectType(payload types.ProjectTypePayload) ([]*types.ProjectType, error) {
	_, err := s.db.Exec(`INSERT INTO project_types (project_type_name, neaktor_id) VALUES ($1, NULLIF($2, ''))`, payload.ProjectTypeName, payload.NeaktorID)
	if err != nil {
		return nil, err
	}
	return s.ListProjectTypes()
}

func (s *Store) UpdateProjectType(id int, payload types.ProjectTypePayload) ([]*types.ProjectType, error) {
	result, err := s.db.Exec(`UPDATE project_types SET project_type_name = $1, neaktor_id = NULLIF($2, '') WHERE project_type_id = $3`, payload.ProjectTypeName, payload.NeaktorID, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListProjectTypes()
}

func (s *Store) DeleteProjectType(id int) ([]*types.ProjectType, error) {
	result, err := s.db.Exec(`DELETE FROM project_types WHERE project_type_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListProjectTypes()
}

func (s *Store) ListWarehouses() ([]*types.Warehouse, error) {
	rows, err := s.db.Query(`SELECT warehouse_id, warehouse_name, COALESCE(warehouse_adress, '') FROM warehouses ORDER BY warehouse_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.Warehouse, 0)
	for rows.Next() {
		item := new(types.Warehouse)
		if err := rows.Scan(&item.WarehouseID, &item.WarehouseName, &item.WarehouseAdress); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) SearchWarehouses(query types.ListQuery) ([]*types.Warehouse, int, error) {
	items, err := s.ListWarehouses()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.Warehouse, 0, len(items))
	for _, item := range items {
		if matchesSearch(search, strconv.Itoa(item.WarehouseID), item.WarehouseName, item.WarehouseAdress) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) CreateWarehouse(payload types.WarehousePayload) ([]*types.Warehouse, error) {
	_, err := s.db.Exec(`INSERT INTO warehouses (warehouse_name, warehouse_adress) VALUES ($1, NULLIF($2, ''))`, payload.WarehouseName, payload.WarehouseAdress)
	if err != nil {
		return nil, err
	}
	return s.ListWarehouses()
}

func (s *Store) UpdateWarehouse(id int, payload types.WarehousePayload) ([]*types.Warehouse, error) {
	result, err := s.db.Exec(`UPDATE warehouses SET warehouse_name = $1, warehouse_adress = NULLIF($2, '') WHERE warehouse_id = $3`, payload.WarehouseName, payload.WarehouseAdress, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListWarehouses()
}

func (s *Store) DeleteWarehouse(id int) ([]*types.Warehouse, error) {
	result, err := s.db.Exec(`DELETE FROM warehouses WHERE warehouse_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListWarehouses()
}

func (s *Store) ListEquipmentSets() ([]*types.EquipmentSet, error) {
	return s.listEquipmentSets("")
}

func (s *Store) SearchEquipmentSets(query types.ListQuery) ([]*types.EquipmentSet, int, error) {
	items, err := s.ListEquipmentSets()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.EquipmentSet, 0, len(items))
	for _, item := range items {
		if matchesSearch(search, strconv.Itoa(item.EquipmentSetID), item.EquipmentSetName, item.Description, item.Type.SetTypeName) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetEquipmentSetByID(id int) (*types.EquipmentSet, error) {
	rows, err := s.listEquipmentSets("WHERE es.equipment_set_id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	return rows[0], nil
}

func (s *Store) CreateEquipmentSet(payload types.EquipmentSetPayload) ([]*types.EquipmentSet, error) {
	setTypeID, err := s.getSetTypeIDByName(payload.SetTypeName)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		INSERT INTO equipment_sets (equipment_set_name, description, set_type_id)
		VALUES ($1, NULLIF($2, ''), $3)
	`, payload.EquipmentSetName, payload.Description, setTypeID)
	if err != nil {
		return nil, err
	}
	return s.ListEquipmentSets()
}

func (s *Store) UpdateEquipmentSet(id int, payload types.EquipmentSetPayload) ([]*types.EquipmentSet, error) {
	setTypeID, err := s.getSetTypeIDByName(payload.SetTypeName)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(`
		UPDATE equipment_sets
		SET equipment_set_name = $1, description = NULLIF($2, ''), set_type_id = $3
		WHERE equipment_set_id = $4
	`, payload.EquipmentSetName, payload.Description, setTypeID, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListEquipmentSets()
}

func (s *Store) DeleteEquipmentSet(id int) ([]*types.EquipmentSet, error) {
	result, err := s.db.Exec(`DELETE FROM equipment_sets WHERE equipment_set_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListEquipmentSets()
}

func (s *Store) GetEquipmentSetsWithMaintenance() ([]*types.EquipmentSet, error) {
	rows, err := s.listEquipmentSets(`
		WHERE EXISTS (
			SELECT 1 FROM equipment e
			WHERE e.equipment_set_id = es.equipment_set_id
			  AND e.needs_maintenance = TRUE
		)
	`)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) GetEquipmentSetsWithStorage() ([]*types.EquipmentSetStorageSummary, error) {
	rows, err := s.db.Query(`
		SELECT
			es.equipment_set_name,
			w.warehouse_name,
			COUNT(e.equipment_id)::INT AS equipment_count
		FROM equipment_sets es
		JOIN equipment e ON e.equipment_set_id = es.equipment_set_id
		JOIN warehouses w ON w.warehouse_id = e.storage_id
		GROUP BY es.equipment_set_name, w.warehouse_name
		ORDER BY es.equipment_set_name, w.warehouse_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.EquipmentSetStorageSummary, 0)
	for rows.Next() {
		item := new(types.EquipmentSetStorageSummary)
		if err := rows.Scan(&item.EquipmentSetName, &item.WarehouseName, &item.EquipmentCount); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) ListEquipment() ([]*types.Equipment, error) {
	return s.listEquipment("")
}

func (s *Store) SearchEquipment(query types.ListQuery) ([]*types.Equipment, int, error) {
	items, err := s.ListEquipment()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.Equipment, 0, len(items))
	for _, item := range items {
		if matchesSearch(
			search,
			strconv.Itoa(item.EquipmentID),
			item.EquipmentName,
			item.SerialNumber,
			item.Description,
			item.EquipmentSet.EquipmentSetName,
			item.Storage.WarehouseName,
		) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) ListEquipmentBySetID(setID int) ([]*types.Equipment, error) {
	return s.listEquipment("WHERE e.equipment_set_id = $1", setID)
}

func (s *Store) SearchEquipmentBySetID(setID int, query types.ListQuery) ([]*types.Equipment, int, error) {
	items, err := s.ListEquipmentBySetID(setID)
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.Equipment, 0, len(items))
	for _, item := range items {
		if matchesSearch(
			search,
			strconv.Itoa(item.EquipmentID),
			item.EquipmentName,
			item.SerialNumber,
			item.Description,
			item.EquipmentSet.EquipmentSetName,
			item.Storage.WarehouseName,
		) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetEquipmentByID(id int) (*types.Equipment, error) {
	rows, err := s.listEquipment("WHERE e.equipment_id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	return rows[0], nil
}

func (s *Store) CreateEquipment(payload types.EquipmentPayload) ([]*types.Equipment, error) {
	equipmentSetID, err := s.getEquipmentSetIDByName(payload.EquipmentSetName)
	if err != nil {
		return nil, err
	}

	warehouseID, err := s.getWarehouseIDByName(payload.WarehouseName)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		INSERT INTO equipment (
			equipment_set_id,
			equipment_name,
			description,
			serial_number,
			storage_id,
			current_storage,
			needs_maintenance,
			date_of_purchase,
			cost_of_purchase
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, NULLIF($6, ''), $7, NULLIF($8, '')::DATE, $9)
	`, equipmentSetID, payload.EquipmentName, payload.Description, payload.SerialNumber, warehouseID, payload.CurrentStorage, payload.NeedsMaintenance, payload.DateOfPurchase, payload.CostOfPurchase)
	if err != nil {
		return nil, err
	}

	return s.ListEquipment()
}

func (s *Store) UpdateEquipment(id int, payload types.EquipmentPayload) ([]*types.Equipment, error) {
	equipmentSetID, err := s.getEquipmentSetIDByName(payload.EquipmentSetName)
	if err != nil {
		return nil, err
	}

	warehouseID, err := s.getWarehouseIDByName(payload.WarehouseName)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(`
		UPDATE equipment
		SET equipment_set_id = $1,
			equipment_name = $2,
			description = NULLIF($3, ''),
			serial_number = $4,
			storage_id = $5,
			current_storage = NULLIF($6, ''),
			needs_maintenance = $7,
			date_of_purchase = NULLIF($8, '')::DATE,
			cost_of_purchase = $9
		WHERE equipment_id = $10
	`, equipmentSetID, payload.EquipmentName, payload.Description, payload.SerialNumber, warehouseID, payload.CurrentStorage, payload.NeedsMaintenance, payload.DateOfPurchase, payload.CostOfPurchase, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}

	return s.ListEquipment()
}

func (s *Store) DeleteEquipment(id int) error {
	result, err := s.db.Exec(`DELETE FROM equipment WHERE equipment_id = $1`, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListProjects(archived bool) ([]*types.Project, error) {
	where := "WHERE p.archived IS FALSE"
	if archived {
		where = "WHERE p.archived IS TRUE"
	}

	result, err := s.listProjects(where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) SearchProjects(archived bool, query types.ListQuery) ([]*types.Project, int, error) {
	items, err := s.ListProjects(archived)
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.Project, 0, len(items))
	for _, item := range items {
		if matchesSearch(
			search,
			strconv.Itoa(item.ProjectID),
			item.ProjectName,
			item.Type.ProjectTypeName,
			item.ChiefEngineer.Name,
			item.ShootingStartDate,
			item.ShootingEndDate,
		) {
			filtered = append(filtered, item)
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetProjectByID(id int) (*types.Project, error) {
	result, err := s.listProjects("WHERE p.project_id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrNotFound
	}

	project := result[0]
	equipment, err := s.listEquipment(`
		WHERE e.equipment_id IN (
			SELECT equipment_id FROM equipment_in_project WHERE project_id = $1
		)
	`, project.ProjectID)
	if err != nil {
		return nil, err
	}
	project.Equipment = equipment
	return project, nil
}

func (s *Store) CreateProject(payload types.ProjectPayload) ([]*types.Project, error) {
	projectTypeID, err := s.getProjectTypeIDByName(payload.ProjectTypeName)
	if err != nil {
		return nil, err
	}

	chiefEngineerID, err := s.getUserIDByName(payload.ChiefEngineerName)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		INSERT INTO projects (
			project_name,
			archived,
			project_type_id,
			shooting_start_date,
			shooting_end_date,
			chief_engineer_id
		)
		VALUES ($1, $2, $3, $4::DATE, $5::DATE, $6)
	`, payload.ProjectName, payload.Archived, projectTypeID, payload.ShootingStartDate, payload.ShootingEndDate, chiefEngineerID)
	if err != nil {
		return nil, err
	}

	return s.ListProjects(false)
}

func (s *Store) UpdateProject(id int, payload types.ProjectPayload) ([]*types.Project, error) {
	projectTypeID, err := s.getProjectTypeIDByName(payload.ProjectTypeName)
	if err != nil {
		return nil, err
	}

	chiefEngineerID, err := s.getUserIDByName(payload.ChiefEngineerName)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(`
		UPDATE projects
		SET project_name = $1,
			archived = $2,
			project_type_id = $3,
			shooting_start_date = $4::DATE,
			shooting_end_date = $5::DATE,
			chief_engineer_id = $6
		WHERE project_id = $7
	`, payload.ProjectName, payload.Archived, projectTypeID, payload.ShootingStartDate, payload.ShootingEndDate, chiefEngineerID, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}

	return s.ListProjects(payload.Archived)
}

func (s *Store) DeleteProject(id int) ([]*types.Project, error) {
	result, err := s.db.Exec(`DELETE FROM projects WHERE project_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListProjects(false)
}

func (s *Store) ListDrafts() ([]*types.Draft, error) {
	rows, err := s.db.Query(`SELECT draft_id, draft_name FROM drafts ORDER BY draft_name DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.Draft, 0)
	for rows.Next() {
		item := new(types.Draft)
		if err := rows.Scan(&item.DraftID, &item.DraftName); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, draft := range result {
		equipment, err := s.listEquipment(`
			WHERE e.equipment_id IN (
				SELECT equipment_id FROM equipment_in_draft WHERE draft_id = $1
			)
		`, draft.DraftID)
		if err != nil {
			return nil, err
		}
		draft.Equipment = equipment
	}

	return result, nil
}

func (s *Store) SearchDrafts(query types.ListQuery) ([]*types.Draft, int, error) {
	items, err := s.ListDrafts()
	if err != nil {
		return nil, 0, err
	}

	search := normalizeSearchQuery(query.Search)
	filtered := make([]*types.Draft, 0, len(items))
	for _, item := range items {
		if matchesSearch(search, strconv.Itoa(item.DraftID), item.DraftName) {
			filtered = append(filtered, item)
			continue
		}

		for _, equipment := range item.Equipment {
			if matchesSearch(search, equipment.EquipmentName, equipment.SerialNumber) {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return paginateSlice(filtered, query), len(filtered), nil
}

func (s *Store) GetDraftByID(id int) (*types.Draft, error) {
	row := s.db.QueryRow(`SELECT draft_id, draft_name FROM drafts WHERE draft_id = $1`, id)
	draft := new(types.Draft)
	if err := row.Scan(&draft.DraftID, &draft.DraftName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	equipment, err := s.listEquipment(`
		WHERE e.equipment_id IN (
			SELECT equipment_id FROM equipment_in_draft WHERE draft_id = $1
		)
	`, draft.DraftID)
	if err != nil {
		return nil, err
	}
	draft.Equipment = equipment
	return draft, nil
}

func (s *Store) CreateDraft(payload types.DraftPayload) ([]*types.Draft, error) {
	_, err := s.db.Exec(`INSERT INTO drafts (draft_name) VALUES ($1)`, payload.DraftName)
	if err != nil {
		return nil, err
	}
	return s.ListDrafts()
}

func (s *Store) UpdateDraft(id int, payload types.DraftPayload) ([]*types.Draft, error) {
	result, err := s.db.Exec(`UPDATE drafts SET draft_name = $1 WHERE draft_id = $2`, payload.DraftName, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListDrafts()
}

func (s *Store) DeleteDraft(id int) ([]*types.Draft, error) {
	result, err := s.db.Exec(`DELETE FROM drafts WHERE draft_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, ErrNotFound
	}
	return s.ListDrafts()
}

func (s *Store) GetEquipmentInProject(projectID int) (*types.EquipmentInProjectResponse, error) {
	return s.buildProjectEquipmentResponse(projectID)
}

func (s *Store) AddEquipmentToProject(payload types.EquipmentInProjectPayload) (*types.EquipmentInProjectResponse, error) {
	_, err := s.db.Exec(`
		INSERT INTO equipment_in_project (project_id, equipment_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, payload.ProjectID, payload.EquipmentID)
	if err != nil {
		return nil, err
	}
	return s.buildProjectEquipmentResponse(payload.ProjectID)
}

func (s *Store) RemoveEquipmentFromProject(payload types.ProjectEquipmentDeletePayload) (*types.EquipmentInProjectResponse, error) {
	_, err := s.db.Exec(`DELETE FROM equipment_in_project WHERE project_id = $1 AND equipment_id = $2`, payload.ProjectID, payload.EquipmentID)
	if err != nil {
		return nil, err
	}
	return s.buildProjectEquipmentResponse(payload.ProjectID)
}

func (s *Store) AddSetToProject(payload types.ProjectSetPayload) (*types.EquipmentInProjectResponse, error) {
	_, err := s.db.Exec(`
		INSERT INTO equipment_in_project (project_id, equipment_id)
		SELECT $1, e.equipment_id
		FROM equipment e
		WHERE e.equipment_set_id = $2
		ON CONFLICT DO NOTHING
	`, payload.ProjectID, payload.EquipmentSetID)
	if err != nil {
		return nil, err
	}
	return s.buildProjectEquipmentResponse(payload.ProjectID)
}

func (s *Store) RemoveSetFromProject(payload types.ProjectSetDeletePayload) (*types.EquipmentInProjectResponse, error) {
	setID, err := s.getEquipmentSetIDByName(payload.EquipmentSetName)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec(`
		DELETE FROM equipment_in_project eip
		USING equipment e
		WHERE eip.project_id = $1
		  AND eip.equipment_id = e.equipment_id
		  AND e.equipment_set_id = $2
	`, payload.ProjectID, setID)
	if err != nil {
		return nil, err
	}
	return s.buildProjectEquipmentResponse(payload.ProjectID)
}

func (s *Store) GetAvailableProjectEquipmentInSet(payload types.ProjectSetPayload) ([]*types.Equipment, error) {
	rows, err := s.listEquipment(`
		WHERE e.equipment_set_id = $1
		  AND e.equipment_id NOT IN (
			SELECT equipment_id FROM equipment_in_project WHERE project_id = $2
		)
	`, payload.EquipmentSetID, payload.ProjectID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) GetConflictingEquipment(projectID int) ([]*types.EquipmentConflict, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT
			e.equipment_id,
			e.equipment_name,
			es.equipment_set_name,
			p2.project_id,
			p2.project_name
		FROM projects p
		JOIN equipment_in_project eip1 ON eip1.project_id = p.project_id
		JOIN equipment_in_project eip2 ON eip2.equipment_id = eip1.equipment_id
		JOIN projects p2 ON p2.project_id = eip2.project_id
		JOIN equipment e ON e.equipment_id = eip1.equipment_id
		JOIN equipment_sets es ON es.equipment_set_id = e.equipment_set_id
		WHERE p.project_id = $1
		  AND p2.project_id <> p.project_id
		  AND p2.archived = FALSE
		  AND NOT (p2.shooting_end_date < p.shooting_start_date OR p2.shooting_start_date > p.shooting_end_date)
		ORDER BY e.equipment_name, p2.project_name
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.EquipmentConflict, 0)
	for rows.Next() {
		item := new(types.EquipmentConflict)
		if err := rows.Scan(&item.EquipmentID, &item.EquipmentName, &item.EquipmentSetName, &item.ProjectID, &item.ProjectName); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) AddDraftToProject(payload types.AddDraftToProjectPayload) (*types.EquipmentInProjectResponse, error) {
	_, err := s.db.Exec(`
		INSERT INTO equipment_in_project (project_id, equipment_id)
		SELECT $1, eid.equipment_id
		FROM equipment_in_draft eid
		WHERE eid.draft_id = $2
		ON CONFLICT DO NOTHING
	`, payload.ProjectID, payload.DraftID)
	if err != nil {
		return nil, err
	}
	return s.buildProjectEquipmentResponse(payload.ProjectID)
}

func (s *Store) ResetEquipmentInProject(projectID int) error {
	_, err := s.db.Exec(`DELETE FROM equipment_in_project WHERE project_id = $1`, projectID)
	return err
}

func (s *Store) GetConflictingProjects() ([]*types.ConflictingProject, error) {
	rows, err := s.db.Query(`
		SELECT
			p.project_id,
			p.project_name,
			TO_CHAR(p.shooting_start_date, 'YYYY-MM-DD') AS shooting_start_date,
			TO_CHAR(p.shooting_end_date, 'YYYY-MM-DD') AS shooting_end_date,
			COUNT(DISTINCT eip1.equipment_id)::INT AS conflicting_equipment_count
		FROM projects p
		JOIN equipment_in_project eip1 ON eip1.project_id = p.project_id
		JOIN equipment_in_project eip2 ON eip2.equipment_id = eip1.equipment_id
		JOIN projects p2 ON p2.project_id = eip2.project_id
		WHERE p.archived = FALSE
		  AND p2.project_id <> p.project_id
		  AND p2.archived = FALSE
		  AND NOT (p2.shooting_end_date < p.shooting_start_date OR p2.shooting_start_date > p.shooting_end_date)
		GROUP BY p.project_id, p.project_name, p.shooting_start_date, p.shooting_end_date
		ORDER BY p.shooting_start_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.ConflictingProject, 0)
	for rows.Next() {
		item := new(types.ConflictingProject)
		if err := rows.Scan(&item.ProjectID, &item.ProjectName, &item.ShootingStartDate, &item.ShootingEndDate, &item.ConflictingEquipmentCount); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) GetEquipmentInDraft(draftID int) (*types.EquipmentInDraftResponse, error) {
	return s.buildDraftEquipmentResponse(draftID)
}

func (s *Store) AddEquipmentToDraft(payload types.EquipmentInDraftPayload) (*types.EquipmentInDraftResponse, error) {
	_, err := s.db.Exec(`
		INSERT INTO equipment_in_draft (draft_id, equipment_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, payload.DraftID, payload.EquipmentID)
	if err != nil {
		return nil, err
	}
	return s.buildDraftEquipmentResponse(payload.DraftID)
}

func (s *Store) RemoveEquipmentFromDraft(payload types.DraftEquipmentDeletePayload) (*types.EquipmentInDraftResponse, error) {
	_, err := s.db.Exec(`DELETE FROM equipment_in_draft WHERE draft_id = $1 AND equipment_id = $2`, payload.DraftID, payload.EquipmentID)
	if err != nil {
		return nil, err
	}
	return s.buildDraftEquipmentResponse(payload.DraftID)
}

func (s *Store) AddSetToDraft(payload types.DraftSetPayload) (*types.EquipmentInDraftResponse, error) {
	_, err := s.db.Exec(`
		INSERT INTO equipment_in_draft (draft_id, equipment_id)
		SELECT $1, e.equipment_id
		FROM equipment e
		WHERE e.equipment_set_id = $2
		ON CONFLICT DO NOTHING
	`, payload.DraftID, payload.EquipmentSetID)
	if err != nil {
		return nil, err
	}
	return s.buildDraftEquipmentResponse(payload.DraftID)
}

func (s *Store) RemoveSetFromDraft(payload types.DraftSetDeletePayload) (*types.EquipmentInDraftResponse, error) {
	setID, err := s.getEquipmentSetIDByName(payload.EquipmentSetName)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec(`
		DELETE FROM equipment_in_draft eid
		USING equipment e
		WHERE eid.draft_id = $1
		  AND eid.equipment_id = e.equipment_id
		  AND e.equipment_set_id = $2
	`, payload.DraftID, setID)
	if err != nil {
		return nil, err
	}
	return s.buildDraftEquipmentResponse(payload.DraftID)
}

func (s *Store) GetAvailableDraftEquipmentInSet(payload types.DraftSetPayload) ([]*types.Equipment, error) {
	rows, err := s.listEquipment(`
		WHERE e.equipment_set_id = $1
		  AND e.equipment_id NOT IN (
			SELECT equipment_id FROM equipment_in_draft WHERE draft_id = $2
		)
	`, payload.EquipmentSetID, payload.DraftID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) buildProjectEquipmentResponse(projectID int) (*types.EquipmentInProjectResponse, error) {
	project, err := s.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}

	equipmentInProject, err := s.listEquipment(`
		WHERE e.equipment_id IN (
			SELECT equipment_id FROM equipment_in_project WHERE project_id = $1
		)
	`, projectID)
	if err != nil {
		return nil, err
	}

	availableEquipment, err := s.listEquipment(`
		WHERE e.equipment_id NOT IN (
			SELECT equipment_id FROM equipment_in_project WHERE project_id = $1
		)
	`, projectID)
	if err != nil {
		return nil, err
	}

	setIDs := uniqueEquipmentSetIDs(availableEquipment)
	setsInProject := make([]*types.EquipmentSet, 0)
	for _, setID := range setIDs {
		set, err := s.GetEquipmentSetByID(setID)
		if err != nil {
			return nil, err
		}
		setsInProject = append(setsInProject, set)
	}

	return &types.EquipmentInProjectResponse{
		Project:            project,
		EquipmentInProject: equipmentInProject,
		AvailableEquipment: availableEquipment,
		SetsInProject:      setsInProject,
	}, nil
}

func (s *Store) buildDraftEquipmentResponse(draftID int) (*types.EquipmentInDraftResponse, error) {
	draft, err := s.GetDraftByID(draftID)
	if err != nil {
		return nil, err
	}

	equipmentInDraft, err := s.listEquipment(`
		WHERE e.equipment_id IN (
			SELECT equipment_id FROM equipment_in_draft WHERE draft_id = $1
		)
	`, draftID)
	if err != nil {
		return nil, err
	}

	availableEquipment, err := s.listEquipment(`
		WHERE e.equipment_id NOT IN (
			SELECT equipment_id FROM equipment_in_draft WHERE draft_id = $1
		)
	`, draftID)
	if err != nil {
		return nil, err
	}

	setIDs := uniqueEquipmentSetIDs(availableEquipment)
	setsInDraft := make([]*types.EquipmentSet, 0)
	for _, setID := range setIDs {
		set, err := s.GetEquipmentSetByID(setID)
		if err != nil {
			return nil, err
		}
		setsInDraft = append(setsInDraft, set)
	}

	return &types.EquipmentInDraftResponse{
		Draft:              draft,
		EquipmentInDraft:   equipmentInDraft,
		AvailableEquipment: availableEquipment,
		SetsInDraft:        setsInDraft,
	}, nil
}

func uniqueEquipmentSetIDs(items []*types.Equipment) []int {
	seen := map[int]struct{}{}
	ids := make([]int, 0)
	for _, item := range items {
		if _, ok := seen[item.EquipmentSetID]; ok {
			continue
		}
		seen[item.EquipmentSetID] = struct{}{}
		ids = append(ids, item.EquipmentSetID)
	}
	return ids
}

func (s *Store) listEquipmentSets(extraWhere string, args ...any) ([]*types.EquipmentSet, error) {
	query := `
		SELECT
			es.equipment_set_id,
			es.equipment_set_name,
			COALESCE(es.description, ''),
			es.set_type_id,
			st.set_type_name
		FROM equipment_sets es
		JOIN set_types st ON st.set_type_id = es.set_type_id
	`
	if strings.TrimSpace(extraWhere) != "" {
		query += " " + extraWhere
	}
	query += " ORDER BY es.equipment_set_name ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.EquipmentSet, 0)
	for rows.Next() {
		item := new(types.EquipmentSet)
		setTypeName := ""
		if err := rows.Scan(&item.EquipmentSetID, &item.EquipmentSetName, &item.Description, &item.SetTypeID, &setTypeName); err != nil {
			return nil, err
		}
		item.Type = &types.SetType{SetTypeID: item.SetTypeID, SetTypeName: setTypeName}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) listEquipment(extraWhere string, args ...any) ([]*types.Equipment, error) {
	query := `
		SELECT
			e.equipment_id,
			e.equipment_set_id,
			e.equipment_name,
			COALESCE(e.description, ''),
			e.serial_number,
			e.storage_id,
			COALESCE(e.current_storage, ''),
			e.needs_maintenance,
			COALESCE(TO_CHAR(e.date_of_purchase, 'YYYY-MM-DD'), ''),
			e.cost_of_purchase,
			es.equipment_set_name,
			COALESCE(es.description, ''),
			es.set_type_id,
			st.set_type_name,
			w.warehouse_name,
			COALESCE(w.warehouse_adress, '')
		FROM equipment e
		JOIN equipment_sets es ON es.equipment_set_id = e.equipment_set_id
		JOIN set_types st ON st.set_type_id = es.set_type_id
		JOIN warehouses w ON w.warehouse_id = e.storage_id
	`
	if strings.TrimSpace(extraWhere) != "" {
		query += " " + extraWhere
	}
	query += " ORDER BY e.equipment_name ASC, e.equipment_id ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.Equipment, 0)
	ids := make([]int, 0)
	for rows.Next() {
		item := new(types.Equipment)
		var cost sql.NullFloat64
		equipmentSetName := ""
		setDescription := ""
		setTypeID := 0
		setTypeName := ""
		warehouseName := ""
		warehouseAdress := ""
		dateOfPurchase := ""

		if err := rows.Scan(
			&item.EquipmentID,
			&item.EquipmentSetID,
			&item.EquipmentName,
			&item.Description,
			&item.SerialNumber,
			&item.StorageID,
			&item.CurrentStorage,
			&item.NeedsMaintenance,
			&dateOfPurchase,
			&cost,
			&equipmentSetName,
			&setDescription,
			&setTypeID,
			&setTypeName,
			&warehouseName,
			&warehouseAdress,
		); err != nil {
			return nil, err
		}

		item.DateOfPurchase = dateOfPurchase
		if cost.Valid {
			item.CostOfPurchase = &cost.Float64
		}
		item.EquipmentSet = &types.EquipmentSet{
			EquipmentSetID:   item.EquipmentSetID,
			EquipmentSetName: equipmentSetName,
			Description:      setDescription,
			SetTypeID:        setTypeID,
			Type: &types.SetType{
				SetTypeID:   setTypeID,
				SetTypeName: setTypeName,
			},
		}
		item.Storage = &types.Warehouse{
			WarehouseID:     item.StorageID,
			WarehouseName:   warehouseName,
			WarehouseAdress: warehouseAdress,
		}
		result = append(result, item)
		ids = append(ids, item.EquipmentID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return result, nil
	}

	projRows, err := s.db.Query(`
		SELECT p.project_id, p.project_name, eip.equipment_id
		FROM equipment_in_project eip
		JOIN projects p ON p.project_id = eip.project_id
		ORDER BY p.project_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer projRows.Close()

	projectMap := map[int][]*types.Project{}
	for projRows.Next() {
		var projectID int
		var projectName string
		var equipmentID int
		if err := projRows.Scan(&projectID, &projectName, &equipmentID); err != nil {
			return nil, err
		}
		projectMap[equipmentID] = append(projectMap[equipmentID], &types.Project{ProjectID: projectID, ProjectName: projectName})
	}
	if err := projRows.Err(); err != nil {
		return nil, err
	}

	for _, item := range result {
		item.Projects = projectMap[item.EquipmentID]
	}

	return result, nil
}

func (s *Store) listProjects(extraWhere string, args ...any) ([]*types.Project, error) {
	query := `
		SELECT
			p.project_id,
			COALESCE(p.neaktor_id, ''),
			p.project_name,
			p.archived,
			COALESCE(p.project_type_id, 0),
			TO_CHAR(p.shooting_start_date, 'YYYY-MM-DD'),
			TO_CHAR(p.shooting_end_date, 'YYYY-MM-DD'),
			COALESCE(p.chief_engineer_id, 0),
			COALESCE(pt.project_type_name, ''),
			COALESCE(u.name, '')
		FROM projects p
		LEFT JOIN project_types pt ON pt.project_type_id = p.project_type_id
		LEFT JOIN users u ON u.id = p.chief_engineer_id
	`
	if strings.TrimSpace(extraWhere) != "" {
		query += " " + extraWhere
	}
	query += " ORDER BY p.shooting_start_date ASC, p.project_id ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*types.Project, 0)
	for rows.Next() {
		item := new(types.Project)
		projectTypeName := ""
		chiefEngineerName := ""
		if err := rows.Scan(
			&item.ProjectID,
			&item.NeaktorID,
			&item.ProjectName,
			&item.Archived,
			&item.ProjectTypeID,
			&item.ShootingStartDate,
			&item.ShootingEndDate,
			&item.ChiefEngineerID,
			&projectTypeName,
			&chiefEngineerName,
		); err != nil {
			return nil, err
		}
		item.Type = &types.ProjectType{ProjectTypeID: item.ProjectTypeID, ProjectTypeName: projectTypeName}
		item.ChiefEngineer = &types.UserShort{ID: item.ChiefEngineerID, Name: chiefEngineerName}
		item.Equipment = []*types.Equipment{}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	equipmentRows, err := s.db.Query(`
		SELECT project_id, equipment_id
		FROM equipment_in_project
		ORDER BY project_id ASC, equipment_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer equipmentRows.Close()

	byProjectID := map[int][]*types.Equipment{}
	for equipmentRows.Next() {
		var projectID int
		var equipmentID int
		if err := equipmentRows.Scan(&projectID, &equipmentID); err != nil {
			return nil, err
		}
		byProjectID[projectID] = append(byProjectID[projectID], &types.Equipment{EquipmentID: equipmentID})
	}
	if err := equipmentRows.Err(); err != nil {
		return nil, err
	}

	for _, project := range result {
		project.Equipment = byProjectID[project.ProjectID]
	}

	return result, nil
}

func (s *Store) getSetTypeIDByName(name string) (int, error) {
	row := s.db.QueryRow(`SELECT set_type_id FROM set_types WHERE set_type_name = $1`, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidReference
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) getProjectTypeIDByName(name string) (int, error) {
	row := s.db.QueryRow(`SELECT project_type_id FROM project_types WHERE project_type_name = $1`, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidReference
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) getWarehouseIDByName(name string) (int, error) {
	row := s.db.QueryRow(`SELECT warehouse_id FROM warehouses WHERE warehouse_name = $1`, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidReference
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) getEquipmentSetIDByName(name string) (int, error) {
	row := s.db.QueryRow(`SELECT equipment_set_id FROM equipment_sets WHERE equipment_set_name = $1`, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidReference
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) getUserIDByName(name string) (int, error) {
	row := s.db.QueryRow(`SELECT id FROM users WHERE name = $1`, name)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidReference
		}
		return 0, err
	}
	return id, nil
}

func mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrInvalidReference) {
		return err
	}
	return fmt.Errorf("store error: %w", err)
}

func normalizeSearchQuery(search string) string {
	return strings.ToLower(strings.TrimSpace(search))
}

func matchesSearch(search string, values ...string) bool {
	if search == "" {
		return true
	}

	for _, value := range values {
		if strings.Contains(strings.ToLower(strings.TrimSpace(value)), search) {
			return true
		}
	}

	return false
}

func paginateSlice[T any](items []T, query types.ListQuery) []T {
	page := query.Page
	if page < 1 {
		page = 1
	}

	perPage := query.PerPage
	if perPage <= 0 {
		perPage = 10
	}

	total := len(items)
	if total == 0 {
		return make([]T, 0)
	}

	start := (page - 1) * perPage
	if start >= total {
		return make([]T, 0)
	}

	end := start + perPage
	if end > total {
		end = total
	}

	return items[start:end]
}
