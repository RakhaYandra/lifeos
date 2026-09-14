-- +goose Up
CREATE TABLE settings (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  active_year INTEGER NOT NULL DEFAULT 2026,
  currency TEXT NOT NULL DEFAULT 'IDR',
  budget_warn_pct INTEGER NOT NULL DEFAULT 80,
  goal_warn_pct INTEGER NOT NULL DEFAULT 70
);

-- +goose Down
DROP TABLE settings;
