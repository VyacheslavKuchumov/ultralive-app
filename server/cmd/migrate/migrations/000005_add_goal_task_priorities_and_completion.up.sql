CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_equipment_set_id ON equipment(equipment_set_id);
CREATE INDEX IF NOT EXISTS idx_equipment_storage_id ON equipment(storage_id);
CREATE INDEX IF NOT EXISTS idx_projects_archived_dates ON projects(archived, shooting_start_date, shooting_end_date);
CREATE INDEX IF NOT EXISTS idx_projects_type_id ON projects(project_type_id);
CREATE INDEX IF NOT EXISTS idx_projects_chief_engineer_id ON projects(chief_engineer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_in_project_equipment_id ON equipment_in_project(equipment_id);
CREATE INDEX IF NOT EXISTS idx_equipment_in_draft_equipment_id ON equipment_in_draft(equipment_id);
