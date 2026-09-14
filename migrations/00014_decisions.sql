-- +goose Up
CREATE TABLE decisions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE TABLE decision_options (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  decision_id INTEGER NOT NULL REFERENCES decisions(id) ON DELETE CASCADE,
  name TEXT NOT NULL
);
CREATE TABLE decision_marks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  option_id INTEGER NOT NULL REFERENCES decision_options(id) ON DELETE CASCADE,
  criterion TEXT NOT NULL,
  weight REAL NOT NULL DEFAULT 1,
  score REAL NOT NULL DEFAULT 0,
  UNIQUE (option_id, criterion)
);
CREATE INDEX idx_options_decision ON decision_options(decision_id);

-- +goose Down
DROP TABLE decision_marks;
DROP TABLE decision_options;
DROP TABLE decisions;
