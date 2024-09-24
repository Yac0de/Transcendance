-- Création de tables et insertion de données de test
CREATE TABLE IF NOT EXISTS test (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL
);

INSERT INTO users (username, email, password)
VALUES
    ('testuser', 'test@example.com', 'password123'),
    ('anotheruser', 'another@example.com', 'password456');
