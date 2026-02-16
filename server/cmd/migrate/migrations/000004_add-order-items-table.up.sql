CREATE TABLE IF NOT EXISTS projects (
  project_id BIGSERIAL PRIMARY KEY,
  neaktor_id TEXT UNIQUE,
  project_name TEXT NOT NULL,
  archived BOOLEAN NOT NULL DEFAULT FALSE,
  project_type_id BIGINT NOT NULL REFERENCES project_types(project_type_id) ON DELETE RESTRICT,
  shooting_start_date DATE NOT NULL,
  shooting_end_date DATE NOT NULL,
  chief_engineer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS drafts (
  draft_id BIGSERIAL PRIMARY KEY,
  draft_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS equipment_in_project (
  project_id BIGINT NOT NULL REFERENCES projects(project_id) ON DELETE CASCADE,
  equipment_id BIGINT NOT NULL REFERENCES equipment(equipment_id) ON DELETE CASCADE,
  PRIMARY KEY (project_id, equipment_id)
);

CREATE TABLE IF NOT EXISTS equipment_in_draft (
  draft_id BIGINT NOT NULL REFERENCES drafts(draft_id) ON DELETE CASCADE,
  equipment_id BIGINT NOT NULL REFERENCES equipment(equipment_id) ON DELETE CASCADE,
  PRIMARY KEY (draft_id, equipment_id)
);
