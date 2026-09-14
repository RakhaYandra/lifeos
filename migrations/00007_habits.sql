-- +goose Up
CREATE TABLE habits (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  frequency TEXT NOT NULL DEFAULT 'daily' CHECK (frequency IN ('daily','weekly')),
  target_per_week INTEGER NOT NULL DEFAULT 1,
  start_date TEXT,
  active INTEGER NOT NULL DEFAULT 1
);
CREATE TABLE habit_logs (
  habit_id INTEGER NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  done INTEGER NOT NULL DEFAULT 1 CHECK (done IN (0,1)),
  PRIMARY KEY (habit_id, date)
);
CREATE INDEX idx_habits_user ON habits(user_id);
CREATE INDEX idx_habit_logs_date ON habit_logs(date);

-- +goose Down
DROP TABLE habit_logs;
DROP TABLE habits;
