CREATE TABLE message_queue (
    msg_uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_body TEXT NOT NULL,
    from_client_id UUID NOT NULL REFERENCES clients(uid),
    to_client_lead VARCHAR(20) NOT NULL,
    scheduled_send_time TIMESTAMP,
    time_sent TIMESTAMP,
    status VARCHAR(50),
    archived BOOLEAN DEFAULT FALSE
);