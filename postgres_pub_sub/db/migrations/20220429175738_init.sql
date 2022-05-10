-- migrate:up
CREATE TYPE status AS ENUM ('available', 'processing', 'error', 'completed', 'deleted');

CREATE TABLE job
(
    id          BIGSERIAL PRIMARY KEY,
    status      status      NOT NULL DEFAULT 'available',
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ,
    delete_time TIMESTAMPTZ,
    data        JSONB
);

CREATE TABLE notification
(
    id          BIGSERIAL PRIMARY KEY,
    job_id      BIGSERIAL REFERENCES job (id),
    status      status,
    prev_status status,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    data        JSONB
);

CREATE OR REPLACE FUNCTION job_notify()
    RETURNS trigger AS
$$
DECLARE
    job_id          BIGINT;
    notification_id BIGINT;
    job_status      status;

BEGIN
    -- Only notify on insert, delete, or update of status
    IF (TG_OP = 'DELETE' OR TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND OLD.status != NEW.status)) THEN
        IF TG_OP = ' INSERT' THEN
            job_id = new.id;
            job_status = new.status;
            INSERT INTO notification (job_id, status, data)
            VALUES (new.id, new.status, new.data)
            RETURNING id INTO notification_id;
        ElSIF TG_OP = ' DELETE' THEN
            job_id = old.id;
            job_status = 'deleted';
            INSERT INTO notification (job_id, status, prev_status, data)
            VALUES (old.id, 'deleted', old.status, old.data)
            RETURNING id INTO notification_id;
        ELSE
            job_id = new.id;
            job_status = new.status;
            INSERT INTO notification (job_id, status, prev_status, data)
            VALUES (new.id, new.status, old.status, new.data)
            RETURNING id INTO notification_id;
        END IF;

        PERFORM pg_notify('job_channel',
                          json_build_object('job_id', job_id, 'notification_id', notification_id, 'status',
                                            job_status)::TEXT);
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER job_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON job
    FOR EACH ROW
EXECUTE PROCEDURE job_notify();

CREATE OR REPLACE FUNCTION job_set_update_time()
    RETURNS trigger AS
$$
BEGIN
    NEW.update_time = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER job_update_trigger
    BEFORE UPDATE
    ON job
    FOR EACH ROW
EXECUTE PROCEDURE job_set_update_time();


-- migrate:down
DROP TRIGGER job_update_trigger ON job;
DROP FUNCTION job_set_update_time;
DROP TRIGGER job_trigger ON job;
DROP FUNCTION job_notify;
DROP TABLE notification;
DROP TABLE job;
DROP TYPE status;