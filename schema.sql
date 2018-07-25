--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.13
-- Dumped by pg_dump version 9.5.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: auction_files; Type: TABLE; Schema: public; Owner: brandon
--

CREATE TABLE public.auction_files (
    id integer NOT NULL,
    url character varying,
    last_modified bigint
);


ALTER TABLE public.auction_files OWNER TO brandon;

--
-- Name: auction_files_id_seq; Type: SEQUENCE; Schema: public; Owner: brandon
--

CREATE SEQUENCE public.auction_files_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auction_files_id_seq OWNER TO brandon;

--
-- Name: auction_files_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: brandon
--

ALTER SEQUENCE public.auction_files_id_seq OWNED BY public.auction_files.id;


--
-- Name: auctions; Type: TABLE; Schema: public; Owner: brandon
--

CREATE TABLE public.auctions (
    id integer NOT NULL,
    auction_id integer,
    item_id integer,
    owner character varying,
    owner_realm character varying,
    bid bigint,
    buyout bigint,
    quantity integer,
    time_left character varying,
    rand integer,
    seed bigint,
    context integer
);


ALTER TABLE public.auctions OWNER TO brandon;

--
-- Name: auctions_id_seq; Type: SEQUENCE; Schema: public; Owner: brandon
--

CREATE SEQUENCE public.auctions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.auctions_id_seq OWNER TO brandon;

--
-- Name: auctions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: brandon
--

ALTER SEQUENCE public.auctions_id_seq OWNED BY public.auctions.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: brandon
--

ALTER TABLE ONLY public.auction_files ALTER COLUMN id SET DEFAULT nextval('public.auction_files_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: brandon
--

ALTER TABLE ONLY public.auctions ALTER COLUMN id SET DEFAULT nextval('public.auctions_id_seq'::regclass);


--
-- Name: auction_files_pkey; Type: CONSTRAINT; Schema: public; Owner: brandon
--

ALTER TABLE ONLY public.auction_files
    ADD CONSTRAINT auction_files_pkey PRIMARY KEY (id);


--
-- Name: auctions_pkey; Type: CONSTRAINT; Schema: public; Owner: brandon
--

ALTER TABLE ONLY public.auctions
    ADD CONSTRAINT auctions_pkey PRIMARY KEY (id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

