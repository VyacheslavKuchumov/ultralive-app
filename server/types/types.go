package types

import "time"

type ListQuery struct {
	Search  string
	Page    int
	PerPage int
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	GetUserByName(name string) (*User, error)
	CreateUser(User) error
	UpdateUserProfile(userID int, payload UpdateProfilePayload) (*User, error)
	UpdateUserPassword(userID int, hashedPassword string) error
	ListUsers() ([]*UserLookup, error)
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserProfile struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserLookup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateProfilePayload struct {
	FirstName       string `json:"firstName" validate:"required,min=1,max=255"`
	LastName        string `json:"lastName" validate:"required,min=1,max=255"`
	Email           string `json:"email" validate:"required,email,max=255"`
	CurrentPassword string `json:"currentPassword" validate:"required,min=3,max=130"`
}

type UpdatePasswordPayload struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=3,max=130"`
	NewPassword     string `json:"newPassword" validate:"required,min=3,max=130"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SetType struct {
	SetTypeID   int    `json:"set_type_id"`
	SetTypeName string `json:"set_type_name"`
}

type SetTypePayload struct {
	SetTypeName string `json:"set_type_name" validate:"required,min=1,max=255"`
}

type ProjectType struct {
	ProjectTypeID   int    `json:"project_type_id"`
	ProjectTypeName string `json:"project_type_name"`
	NeaktorID       string `json:"neaktor_id,omitempty"`
}

type ProjectTypePayload struct {
	ProjectTypeName string `json:"project_type_name" validate:"required,min=1,max=255"`
	NeaktorID       string `json:"neaktor_id,omitempty"`
}

type Warehouse struct {
	WarehouseID     int    `json:"warehouse_id"`
	WarehouseName   string `json:"warehouse_name"`
	WarehouseAdress string `json:"warehouse_adress,omitempty"`
}

type WarehousePayload struct {
	WarehouseName   string `json:"warehouse_name" validate:"required,min=1,max=255"`
	WarehouseAdress string `json:"warehouse_adress"`
}

type EquipmentSet struct {
	EquipmentSetID   int      `json:"equipment_set_id"`
	EquipmentSetName string   `json:"equipment_set_name"`
	Description      string   `json:"description,omitempty"`
	SetTypeID        int      `json:"set_type_id"`
	Type             *SetType `json:"type,omitempty"`
}

type EquipmentSetPayload struct {
	EquipmentSetName string `json:"equipment_set_name" validate:"required,min=1,max=255"`
	Description      string `json:"description"`
	SetTypeName      string `json:"set_type_name" validate:"required,min=1,max=255"`
}

type EquipmentSetStorageSummary struct {
	EquipmentSetName string `json:"equipment_set_name"`
	WarehouseName    string `json:"warehouse_name"`
	EquipmentCount   int    `json:"equipment_count"`
}

type Equipment struct {
	EquipmentID      int           `json:"equipment_id"`
	EquipmentSetID   int           `json:"equipment_set_id"`
	EquipmentName    string        `json:"equipment_name"`
	Description      string        `json:"description,omitempty"`
	SerialNumber     string        `json:"serial_number"`
	StorageID        int           `json:"storage_id"`
	CurrentStorage   string        `json:"current_storage,omitempty"`
	NeedsMaintenance bool          `json:"needs_maintenance"`
	DateOfPurchase   string        `json:"date_of_purchase,omitempty"`
	CostOfPurchase   *float64      `json:"cost_of_purchase,omitempty"`
	EquipmentSet     *EquipmentSet `json:"equipment_set,omitempty"`
	Storage          *Warehouse    `json:"storage,omitempty"`
	Projects         []*Project    `json:"projects,omitempty"`
}

type EquipmentPayload struct {
	EquipmentName    string   `json:"equipment_name" validate:"required,min=1,max=255"`
	SerialNumber     string   `json:"serial_number" validate:"required,min=1,max=255"`
	EquipmentSetName string   `json:"equipment_set_name" validate:"required,min=1,max=255"`
	Description      string   `json:"description"`
	WarehouseName    string   `json:"warehouse_name" validate:"required,min=1,max=255"`
	CurrentStorage   string   `json:"current_storage_name"`
	NeedsMaintenance bool     `json:"needs_maintenance"`
	DateOfPurchase   string   `json:"date_of_purchase"`
	CostOfPurchase   *float64 `json:"cost_of_purchase"`
}

type UserShort struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ProjectID         int          `json:"project_id"`
	NeaktorID         string       `json:"neaktor_id,omitempty"`
	ProjectName       string       `json:"project_name"`
	Archived          bool         `json:"archived"`
	ProjectTypeID     int          `json:"project_type_id"`
	ShootingStartDate string       `json:"shooting_start_date"`
	ShootingEndDate   string       `json:"shooting_end_date"`
	ChiefEngineerID   int          `json:"chief_engineer_id"`
	Type              *ProjectType `json:"type,omitempty"`
	ChiefEngineer     *UserShort   `json:"chiefEngineer,omitempty"`
	Equipment         []*Equipment `json:"equipment,omitempty"`
}

type ProjectPayload struct {
	ProjectName       string `json:"project_name" validate:"required,min=1,max=255"`
	ProjectTypeName   string `json:"project_type_name" validate:"required,min=1,max=255"`
	Archived          bool   `json:"archived"`
	ChiefEngineerName string `json:"chief_engineer_name" validate:"required,min=1,max=255"`
	ShootingStartDate string `json:"shooting_start_date" validate:"required,len=10"`
	ShootingEndDate   string `json:"shooting_end_date" validate:"required,len=10"`
}

type Draft struct {
	DraftID   int          `json:"draft_id"`
	DraftName string       `json:"draft_name"`
	Equipment []*Equipment `json:"equipment,omitempty"`
}

type DraftPayload struct {
	DraftName string `json:"draft_name" validate:"required,min=1,max=255"`
}

type EquipmentInProjectResponse struct {
	Project            *Project        `json:"project"`
	EquipmentInProject []*Equipment    `json:"equipment_in_project"`
	AvailableEquipment []*Equipment    `json:"available_equipment"`
	SetsInProject      []*EquipmentSet `json:"sets_in_project"`
}

type EquipmentInDraftResponse struct {
	Draft              *Draft          `json:"draft"`
	EquipmentInDraft   []*Equipment    `json:"equipment_in_draft"`
	AvailableEquipment []*Equipment    `json:"available_equipment"`
	SetsInDraft        []*EquipmentSet `json:"sets_in_draft"`
}

type EquipmentInProjectPayload struct {
	ProjectID      int `json:"project_id" validate:"required,min=1"`
	EquipmentID    int `json:"equipment_id" validate:"required,min=1"`
	EquipmentSetID int `json:"equipment_set_id,omitempty"`
}

type ProjectEquipmentDeletePayload struct {
	ProjectID   int `json:"project_id" validate:"required,min=1"`
	EquipmentID int `json:"equipment_id" validate:"required,min=1"`
}

type ProjectSetPayload struct {
	ProjectID      int `json:"project_id" validate:"required,min=1"`
	EquipmentSetID int `json:"equipment_set_id" validate:"required,min=1"`
}

type ProjectSetDeletePayload struct {
	ProjectID        int    `json:"project_id" validate:"required,min=1"`
	EquipmentSetName string `json:"equipment_set_name" validate:"required,min=1,max=255"`
}

type AddDraftToProjectPayload struct {
	ProjectID int `json:"project_id" validate:"required,min=1"`
	DraftID   int `json:"draft_id" validate:"required,min=1"`
}

type EquipmentInDraftPayload struct {
	DraftID        int `json:"draft_id" validate:"required,min=1"`
	EquipmentID    int `json:"equipment_id" validate:"required,min=1"`
	EquipmentSetID int `json:"equipment_set_id,omitempty"`
}

type DraftEquipmentDeletePayload struct {
	DraftID     int `json:"draft_id" validate:"required,min=1"`
	EquipmentID int `json:"equipment_id" validate:"required,min=1"`
}

type DraftSetPayload struct {
	DraftID        int `json:"draft_id" validate:"required,min=1"`
	EquipmentSetID int `json:"equipment_set_id" validate:"required,min=1"`
}

type DraftSetDeletePayload struct {
	DraftID          int    `json:"draft_id" validate:"required,min=1"`
	EquipmentSetName string `json:"equipment_set_name" validate:"required,min=1,max=255"`
}

type EquipmentConflict struct {
	EquipmentID      int    `json:"equipment_id"`
	EquipmentName    string `json:"equipment_name"`
	EquipmentSetName string `json:"equipment_set_name"`
	ProjectID        int    `json:"project_id"`
	ProjectName      string `json:"project_name"`
}

type ConflictingProject struct {
	ProjectID                 int    `json:"project_id"`
	ProjectName               string `json:"project_name"`
	ShootingStartDate         string `json:"shooting_start_date"`
	ShootingEndDate           string `json:"shooting_end_date"`
	ConflictingEquipmentCount int    `json:"conflicting_equipment_count"`
}
