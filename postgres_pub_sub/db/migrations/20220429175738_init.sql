-- migrate:up
CREATE TABLE job
(
    id   SERIAL PRIMARY KEY,
    data JSONB
);

CREATE OR REPLACE FUNCTION job_notify()
    RETURNS trigger AS
$$
BEGIN
    PERFORM pg_notify('job_channel', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER job_trigger
    AFTER INSERT OR UPDATE OF data
    ON job
    FOR EACH ROW
EXECUTE PROCEDURE job_notify();
-- migrate:down

DROP TRIGGER job_trigger ON job;
DROP FUNCTION job_notify;
DROP TABLE job;