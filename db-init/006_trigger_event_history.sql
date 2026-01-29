CREATE OR REPLACE FUNCTION log_message_queue_status_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Only log when status actually changes
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
