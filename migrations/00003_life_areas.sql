-- +goose Up
CREATE TABLE life_areas (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

INSERT INTO life_areas (name) VALUES
  ('Career'),
  ('Finance'),
  ('Health & Fitness'),
  ('Learning'),
  ('Personal Development'),
  ('Relationships'),
  ('Family'),
  ('Home / Living');

-- +goose Down
DROP TABLE life_areas;
