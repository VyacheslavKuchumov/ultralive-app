CREATE TABLE IF NOT EXISTS set_types (
  set_type_id BIGSERIAL PRIMARY KEY,
  set_type_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS project_types (
  project_type_id BIGSERIAL PRIMARY KEY,
  project_type_name TEXT NOT NULL UNIQUE,
  neaktor_id TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS warehouses (
  warehouse_id BIGSERIAL PRIMARY KEY,
  warehouse_name TEXT NOT NULL UNIQUE,
  warehouse_adress TEXT
);
