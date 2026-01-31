-- Trigger function for INSERT: log initial 'queued' status
CREATE OR REPLACE FUNCTION log_message_queue_initial_status()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO message_event_history (
        msg_uid,
        event_time,
        prev_status,
        curr_status
    )
    VALUES (
        NEW.msg_uid,
        NOW(),
        NULL,
        NEW.status
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for UPDATE: log status changes
CREATE OR REPLACE FUNCTION log_message_queue_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IS DISTINCT FROM OLD.status THEN
        INSERT INTO message_event_history (
            msg_uid,
            event_time,
            prev_status,
            curr_status
        )
        VALUES (
            OLD.msg_uid,
            NOW(),
            OLD.status,
            NEW.status
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- $$ LANGUAGE plpgsql;

-- Apply triggers
CREATE TRIGGER log_initial_status
    AFTER INSERT ON message_queue
    FOR EACH ROW
    EXECUTE FUNCTION log_message_queue_initial_status();

CREATE TRIGGER log_status_change
    AFTER UPDATE ON message_queue
    FOR EACH ROW
    EXECUTE FUNCTION log_message_queue_status_change();   
