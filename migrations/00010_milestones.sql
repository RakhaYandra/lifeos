-- +goose Up
CREATE TABLE milestones (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  goal_id INTEGER REFERENCES goals(id) ON DELETE SET NULL,
  project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  target_date TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed','cancelled')),
  completed_at TEXT,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_milestones_goal ON milestones(goal_id);
CREATE INDEX idx_milestones_user_date ON milestones(user_id, target_date);

-- +goose Down
DROP TABLE milestones;
