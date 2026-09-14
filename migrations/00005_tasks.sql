-- +goose Up
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
  goal_id INTEGER,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  status TEXT NOT NULL DEFAULT 'inbox' CHECK (status IN ('inbox','not_started','in_progress','waiting','completed','cancelled')),
  priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low','medium','high','critical')),
  due_date TEXT,
  completed_at TEXT,
  effort_est REAL,
  effort_actual REAL
);
CREATE INDEX idx_tasks_user_due ON tasks(user_id, due_date);
CREATE INDEX idx_tasks_project ON tasks(project_id);

-- +goose Down
DROP TABLE tasks;
