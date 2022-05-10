SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: adapters; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA adapters;


--
-- Name: status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.status AS ENUM (
    'available',
    'processing',
    'error',
    'completed',
    'deleted'
);


--
-- Name: job_notify(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.job_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;


--
-- Name: job_set_update_time(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.job_set_update_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.update_time = now();
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: axs_event_details; Type: TABLE; Schema: adapters; Owner: -
--

CREATE TABLE adapters.axs_event_details (
    id timestamp without time zone,
    venue text,
    ticket text,
    context text,
    axs_event text,
    zone text
);


--
-- Name: axs_ticket_sale_state_change; Type: TABLE; Schema: adapters; Owner: -
--

CREATE TABLE adapters.axs_ticket_sale_state_change (
    id bigint NOT NULL,
    _env text,
    _timestamp timestamp with time zone,
    action integer,
    context_id integer,
    event_id integer,
    zone_id integer,
    seat_status_id integer,
    price_code_id integer,
    seat_id integer,
    seat_state_id integer,
    history_pointer bigint,
    price_level_id integer,
    seat_previous_status_id integer,
    seat_previous_state_id integer,
    outlet_id integer,
    channel_id integer,
    price numeric,
    seat_group_id integer,
    action_name text,
    history_date timestamp with time zone,
    seat_state_name text,
    seat_previous_state_name text,
    order_number integer,
    cart_id integer,
    seat_previous_status_name text,
    seat_status_name text,
    seat_key text,
    stream_event_id text,
    stream_partition_id text
);


--
-- Name: axs_ticket_sale_state_change_id_seq; Type: SEQUENCE; Schema: adapters; Owner: -
--

CREATE SEQUENCE adapters.axs_ticket_sale_state_change_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: axs_ticket_sale_state_change_id_seq; Type: SEQUENCE OWNED BY; Schema: adapters; Owner: -
--

ALTER SEQUENCE adapters.axs_ticket_sale_state_change_id_seq OWNED BY adapters.axs_ticket_sale_state_change.id;


--
-- Name: event_details; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.event_details (
    id timestamp without time zone,
    venue text,
    excluded boolean,
    updated_datetime timestamp with time zone,
    override_updated_at timestamp with time zone,
    override_id integer,
    override_distributed integer,
    name text,
    performer text,
    category text,
    sub_category text,
    ticketing_ids text,
    ticketing_start timestamp without time zone,
    ticketing_end timestamp without time zone,
    concession_ids text,
    concession_start timestamp without time zone,
    concession_end timestamp without time zone,
    incident_ids text,
    incident_start timestamp without time zone,
    incident_end timestamp without time zone,
    merchandising_ids text,
    merchandising_start timestamp without time zone,
    merchandising_end timestamp without time zone,
    satisfaction_ids text,
    satisfaction_start timestamp without time zone,
    satisfaction_end timestamp without time zone,
    parking_ids text,
    parking_start timestamp without time zone,
    parking_end timestamp without time zone,
    owners json
);


--
-- Name: job; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job (
    id bigint NOT NULL,
    status public.status DEFAULT 'available'::public.status NOT NULL,
    create_time timestamp with time zone DEFAULT now() NOT NULL,
    update_time timestamp with time zone,
    delete_time timestamp with time zone,
    data jsonb
);


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.job_id_seq OWNED BY public.job.id;


--
-- Name: notification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notification (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    status public.status,
    prev_status public.status,
    create_time timestamp with time zone DEFAULT now() NOT NULL,
    data jsonb
);


--
-- Name: notification_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.notification_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: notification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.notification_id_seq OWNED BY public.notification.id;


--
-- Name: notification_job_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.notification_job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: notification_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.notification_job_id_seq OWNED BY public.notification.job_id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: axs_ticket_sale_state_change id; Type: DEFAULT; Schema: adapters; Owner: -
--

ALTER TABLE ONLY adapters.axs_ticket_sale_state_change ALTER COLUMN id SET DEFAULT nextval('adapters.axs_ticket_sale_state_change_id_seq'::regclass);


--
-- Name: job id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job ALTER COLUMN id SET DEFAULT nextval('public.job_id_seq'::regclass);


--
-- Name: notification id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notification ALTER COLUMN id SET DEFAULT nextval('public.notification_id_seq'::regclass);


--
-- Name: notification job_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notification ALTER COLUMN job_id SET DEFAULT nextval('public.notification_job_id_seq'::regclass);


--
-- Name: axs_ticket_sale_state_change axs_ticket_sale_state_change_pkey; Type: CONSTRAINT; Schema: adapters; Owner: -
--

ALTER TABLE ONLY adapters.axs_ticket_sale_state_change
    ADD CONSTRAINT axs_ticket_sale_state_change_pkey PRIMARY KEY (id);


--
-- Name: job job_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_pkey PRIMARY KEY (id);


--
-- Name: notification notification_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_axs_ticket_state_change_context; Type: INDEX; Schema: adapters; Owner: -
--

CREATE INDEX idx_axs_ticket_state_change_context ON adapters.axs_ticket_sale_state_change USING btree (context_id, event_id, zone_id, seat_id);


--
-- Name: job job_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER job_trigger AFTER INSERT OR DELETE OR UPDATE ON public.job FOR EACH ROW EXECUTE FUNCTION public.job_notify();


--
-- Name: job job_update_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER job_update_trigger BEFORE UPDATE ON public.job FOR EACH ROW EXECUTE FUNCTION public.job_set_update_time();


--
-- Name: notification notification_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20220424023931'),
    ('20220429175738');
