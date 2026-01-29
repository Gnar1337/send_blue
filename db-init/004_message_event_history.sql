CREATE TABLE message_event_history (
    msg_uid UUID NOT NULL,
    time_stamp TIMESTAMP NOT NULL,
    prev_status VARCHAR(255),
    curr_status VARCHAR(255),
    FOREIGN KEY (msg_uid) REFERENCES message_queue(msg_uid)
);