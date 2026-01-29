CREATE TABLE message_archive (
    msg_uid UUID  REFERENCES message_queue(msg_uid),
        message_body TEXT NOT NULL,
        from_client_id UUID NOT NULL REFERENCES clients(uid),
        to_client_lead VARCHAR(20) NOT NULL REFERENCES client_leads(lead_number),
        time_archived TIMESTAMP,
        status VARCHAR(50)
);