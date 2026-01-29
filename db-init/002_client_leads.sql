CREATE TABLE client_leads (
    lead_number VARCHAR(20) PRIMARY KEY,
    client_uid UUID NOT NULL,
    messages_received INT DEFAULT 0,
    last_contacted TIMESTAMP,
    FOREIGN KEY (client_uid) REFERENCES clients(uid)
);
