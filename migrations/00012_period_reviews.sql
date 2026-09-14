-- +goose Up
CREATE TABLE monthly_reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  period TEXT NOT NULL,
  stats TEXT NOT NULL DEFAULT '{}',
  wins TEXT NOT NULL DEFAULT '',
  challenges TEXT NOT NULL DEFAULT '',
  lessons TEXT NOT NULL DEFAULT '',
  next_focus TEXT NOT NULL DEFAULT '',
  UNIQUE (user_id, period)
);
CREATE TABLE yearly_reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  period TEXT NOT NULL,
  stats TEXT NOT NULL DEFAULT '{}',
  wins TEXT NOT NULL DEFAULT '',
  challenges TEXT NOT NULL DEFAULT '',
  lessons TEXT NOT NULL DEFAULT '',
  next_focus TEXT NOT NULL DEFAULT '',
  achievements TEXT NOT NULL DEFAULT '',
  next_year TEXT NOT NULL DEFAULT '',
  UNIQUE (user_id, period)
);

-- +goose Down
DROP TABLE yearly_reviews;
DROP TABLE monthly_reviews;
