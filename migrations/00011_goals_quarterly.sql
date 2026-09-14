-- +goose Up
CREATE TABLE goals_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  level TEXT NOT NULL CHECK (level IN ('annual','quarterly','monthly')),
  parent_id INTEGER REFERENCES goals_new(id) ON DELETE SET NULL,
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  metric TEXT NOT NULL DEFAULT '',
  target_value REAL NOT NULL DEFAULT 0,
  current_value REAL NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('not_started','active','on_track','at_risk','completed','cancelled')),
  target_date TEXT
);
INSERT INTO goals_new SELECT * FROM goals;
DROP TABLE goals;
ALTER TABLE goals_new RENAME TO goals;
CREATE INDEX idx_goals_user_level ON goals(user_id, level);
CREATE INDEX idx_goals_parent ON goals(parent_id);

-- +goose Down
CREATE TABLE goals_old (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  level TEXT NOT NULL CHECK (level IN ('annual','monthly')),
  parent_id INTEGER REFERENCES goals_old(id) ON DELETE SET NULL,
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  metric TEXT NOT NULL DEFAULT '',
  target_value REAL NOT NULL DEFAULT 0,
  current_value REAL NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('not_started','active','on_track','at_risk','completed','cancelled')),
  target_date TEXT
);
INSERT INTO goals_old SELECT * FROM goals WHERE level != 'quarterly';
DROP TABLE goals;
ALTER TABLE goals_old RENAME TO goals;
CREATE INDEX idx_goals_user_level ON goals(user_id, level);
CREATE INDEX idx_goals_parent ON goals(parent_id);
