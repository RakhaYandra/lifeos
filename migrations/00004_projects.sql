-- +goose Up
CREATE TABLE projects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  goal_id INTEGER,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('planning','active','on_hold','completed','cancelled')),
  start_date TEXT,
  target_date TEXT,
  completed_at TEXT,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_projects_user ON projects(user_id);

-- +goose Down
DROP TABLE projects;
