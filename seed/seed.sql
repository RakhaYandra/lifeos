-- Seed P0 fiktif (bukan data asli). Password: Rahasia123
INSERT INTO users (email, password_hash) VALUES
('aku@lifeos.local', '$2a$10$GMr6NAxuAxPKutSu2J5MUO3DG.3iqRcYv5Gu9NQNMnTNRErbsKqMu');

INSERT INTO settings (user_id, active_year, currency, budget_warn_pct, goal_warn_pct) VALUES
(1, 2026, 'IDR', 80, 70);
