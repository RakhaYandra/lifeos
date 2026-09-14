-- +goose Up
CREATE TABLE transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('income','expense','transfer')),
  category TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  amount REAL NOT NULL CHECK (amount > 0),
  account TEXT NOT NULL DEFAULT '',
  recurring INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_trx_user_date ON transactions(user_id, date);
CREATE INDEX idx_trx_cat ON transactions(user_id, category);

CREATE TABLE budgets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  year INTEGER NOT NULL,
  month INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
  category TEXT NOT NULL,
  amount REAL NOT NULL CHECK (amount > 0),
  UNIQUE (user_id, year, month, category)
);

CREATE TABLE subscriptions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  service TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT '',
  cost REAL NOT NULL CHECK (cost >= 0),
  frequency TEXT NOT NULL DEFAULT 'monthly' CHECK (frequency IN ('weekly','monthly','yearly')),
  next_billing TEXT NOT NULL,
  payment_method TEXT NOT NULL DEFAULT '',
  auto_renew INTEGER NOT NULL DEFAULT 1,
  active INTEGER NOT NULL DEFAULT 1,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX idx_sub_user_bill ON subscriptions(user_id, next_billing);

-- +goose Down
DROP TABLE subscriptions;
DROP TABLE budgets;
DROP TABLE transactions;
