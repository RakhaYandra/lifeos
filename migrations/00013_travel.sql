-- +goose Up
CREATE TABLE trips (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  destination TEXT NOT NULL DEFAULT '',
  start_date TEXT NOT NULL,
  end_date TEXT NOT NULL,
  budget REAL NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'planning' CHECK (status IN ('planning','booked','ongoing','done','cancelled')),
  notes TEXT NOT NULL DEFAULT ''
);
CREATE TABLE itinerary_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  time TEXT NOT NULL DEFAULT '',
  activity TEXT NOT NULL,
  location TEXT NOT NULL DEFAULT '',
  cost REAL NOT NULL DEFAULT 0,
  booked INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE packing_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
  category TEXT NOT NULL DEFAULT '',
  item TEXT NOT NULL,
  qty INTEGER NOT NULL DEFAULT 1,
  packed INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_itinerary_trip ON itinerary_items(trip_id, date);
CREATE INDEX idx_packing_trip ON packing_items(trip_id);

-- +goose Down
DROP TABLE packing_items;
DROP TABLE itinerary_items;
DROP TABLE trips;
