--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4 (Debian 16.4-1.pgdg120+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: test; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.test (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.test OWNER TO myuser;

--
-- Data for Name: test; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.test (id, name) FROM stdin;
0	yanis
\.


--
-- PostgreSQL database dump complete
--

