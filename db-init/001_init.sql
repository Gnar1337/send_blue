CREATE TABLE clients (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL,
        messages_sent INT DEFAULT 0
);
