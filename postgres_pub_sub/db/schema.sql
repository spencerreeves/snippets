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
-- Name: job_notify(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.job_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    PERFORM pg_notify('job_channel', NEW.id::text);
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

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
    history_pointer integer,
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
-- Name: job; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job (
    id integer NOT NULL,
    data jsonb
);


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.job_id_seq
    AS integer
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

CREATE TRIGGER job_trigger AFTER INSERT OR UPDATE OF data ON public.job FOR EACH ROW EXECUTE FUNCTION public.job_notify();


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20220429175738');
