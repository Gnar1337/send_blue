CREATE TABLE clients (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL,
        messages_sent INT DEFAULT 0
);
INSERT INTO clients (name, messages_sent) VALUES
    ('Me', 0),
    ('You', 0),
    ('We', 0),
    ('SendBlue', 0);
