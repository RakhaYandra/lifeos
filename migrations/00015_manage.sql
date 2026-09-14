-- +goose Up
CREATE TABLE savings_goals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  target_amount REAL NOT NULL CHECK (target_amount > 0),
  current_amount REAL NOT NULL DEFAULT 0,
  monthly_contribution REAL NOT NULL DEFAULT 0,
  target_date TEXT,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','paused','completed','cancelled'))
);

CREATE TABLE assets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  purchase_price REAL NOT NULL DEFAULT 0,
  current_value REAL NOT NULL DEFAULT 0,
  condition TEXT NOT NULL DEFAULT 'good' CHECK (condition IN ('new','good','worn','broken')),
  location TEXT NOT NULL DEFAULT '',
  warranty_end TEXT,
  notes TEXT NOT NULL DEFAULT ''
);

CREATE TABLE wishlist (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  item TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low','medium','high')),
  est_price REAL NOT NULL DEFAULT 0,
  saved REAL NOT NULL DEFAULT 0,
  target_date TEXT,
  status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned','saving','bought','dropped'))
);

CREATE TABLE documents (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  item TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  expiry_date TEXT NOT NULL,
  reminder_days INTEGER NOT NULL DEFAULT 30,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_documents_user_expiry ON documents(user_id, expiry_date);

CREATE TABLE contacts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  relation TEXT NOT NULL DEFAULT '',
  birthday TEXT,
  last_contact TEXT,
  contact_method TEXT NOT NULL DEFAULT '',
  followup_days INTEGER NOT NULL DEFAULT 30,
  notes TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE contacts;
DROP TABLE documents;
DROP TABLE wishlist;
DROP TABLE assets;
DROP TABLE savings_goals;
