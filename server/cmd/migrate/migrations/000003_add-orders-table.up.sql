CREATE TABLE IF NOT EXISTS equipment_sets (
  equipment_set_id BIGSERIAL PRIMARY KEY,
  equipment_set_name TEXT NOT NULL UNIQUE,
  description TEXT,
  set_type_id BIGINT NOT NULL REFERENCES set_types(set_type_id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS equipment (
  equipment_id BIGSERIAL PRIMARY KEY,
  equipment_set_id BIGINT NOT NULL REFERENCES equipment_sets(equipment_set_id) ON DELETE RESTRICT,
  equipment_name TEXT NOT NULL,
  description TEXT,
  serial_number TEXT NOT NULL,
  storage_id BIGINT NOT NULL REFERENCES warehouses(warehouse_id) ON DELETE RESTRICT,
  current_storage TEXT,
  needs_maintenance BOOLEAN NOT NULL DEFAULT FALSE,
  date_of_purchase DATE,
  cost_of_purchase NUMERIC(12, 2)
);
