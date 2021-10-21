CREATE ROLE api WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'xxxxxx';

GRANT pg_read_all_data, pg_write_all_data TO api;

CREATE DATABASE test_db
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

USE test_db;
CREATE EXTENSION "uuid-ossp";

CREATE TABLE public.users
(
    id bigserial NOT NULL,
    username text NOT NULL,
    password_hash text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uk_users_username UNIQUE (username)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

CREATE TABLE public.sessions
(
    id bigserial NOT NULL,
    token uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id bigint NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_sessions_users_user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE IF EXISTS public.sessions
    OWNER to postgres;