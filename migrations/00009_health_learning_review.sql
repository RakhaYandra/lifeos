-- +goose Up
CREATE TABLE health_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  weight REAL,
  sleep_hours REAL,
  water_liters REAL,
  energy INTEGER CHECK (energy BETWEEN 1 AND 5),
  mood INTEGER CHECK (mood BETWEEN 1 AND 5),
  notes TEXT NOT NULL DEFAULT '',
  UNIQUE (user_id, date)
);

CREATE TABLE workouts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  type TEXT NOT NULL,
  duration_min INTEGER NOT NULL DEFAULT 0,
  intensity TEXT NOT NULL DEFAULT 'medium' CHECK (intensity IN ('low','medium','high')),
  calories INTEGER,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_workouts_user_date ON workouts(user_id, date);

CREATE TABLE learning_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  topic TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'course' CHECK (type IN ('course','book','tutorial','certification','practice','other')),
  provider TEXT NOT NULL DEFAULT '',
  related_skill TEXT NOT NULL DEFAULT '',
  progress INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  hours REAL NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('not_started','active','completed','cancelled')),
  notes TEXT NOT NULL DEFAULT ''
);

CREATE TABLE reading_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'book' CHECK (type IN ('book','article','research','docs','other')),
  author TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'reading' CHECK (status IN ('planned','reading','finished','dropped')),
  rating INTEGER CHECK (rating BETWEEN 1 AND 5),
  takeaways TEXT NOT NULL DEFAULT ''
);

CREATE TABLE reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  week_start TEXT NOT NULL,
  stats TEXT NOT NULL DEFAULT '{}',
  wins TEXT NOT NULL DEFAULT '',
  challenges TEXT NOT NULL DEFAULT '',
  lessons TEXT NOT NULL DEFAULT '',
  next_focus TEXT NOT NULL DEFAULT '',
  UNIQUE (user_id, week_start)
);

CREATE TABLE reminders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  date TEXT NOT NULL,
  recurrence TEXT NOT NULL DEFAULT 'none' CHECK (recurrence IN ('none','daily','weekly','monthly','yearly')),
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_reminders_user_date ON reminders(user_id, date);

-- +goose Down
DROP TABLE reminders;
DROP TABLE reviews;
DROP TABLE reading_entries;
DROP TABLE learning_entries;
DROP TABLE workouts;
DROP TABLE health_logs;
