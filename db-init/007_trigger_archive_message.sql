CREATE OR REPLACE FUNCTION archive_message_on_received()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert into archive table
    INSERT INTO message_archive (
        msg_uid,
        message_body,
        from_client_id,
        to_client_lead,
        time_archived,
        status
    )
    VALUES (
        NEW.msg_uid,
        NEW.message_body,
        NEW.from_client_id,
        NEW.to_client_lead,
        NOW(),
        NEW.status
    );

    -- Remove from queue
    DELETE FROM message_queue
    WHERE msg_uid = NEW.msg_uid;

    RETURN NULL;  -- AFTER triggers can return NULL
END;
$$ LANGUAGE plpgsql;