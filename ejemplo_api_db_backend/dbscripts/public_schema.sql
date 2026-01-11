-- Initial public schema relates to Library 0.x

-- No definimos ningún timeout para la ejecución de las sentencias, bloqueos o transacciones inactivas
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;

-- usamos el conjunto de caracteres UTF8
SET client_encoding = 'UTF8';
-- definimos la zona horaria por defecto como UTC
SET timezone = 'UTC';
-- definimos el formato de los números para que use el punto como separador decimal
SET numeric_std = 'on';
-- definimos el comportamiento de las comillas simples en las cadenas de texto  
SET standard_conforming_strings = on;
-- nivel de mensajes mínimos a mostrar
SET client_min_messages = warning;
-- desactivamos la seguridad a nivel de fila
SET row_security = off;

-- extensión que permite usar el lenguaje PL/pgSQL en funciones y triggers
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
-- extensión que se utiliza para generar UUIDs. Incluye el tipo uuid que usamos en las columnas id de las tablas runners y results
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

SET search_path = public, pg_catalog;
SET default_tablespace = '';

-- runners
CREATE TABLE runners (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(), -- generamos un UUID basado en la dirección MAC del servidor y la fecha/hora actual
    first_name text NOT NULL,
    last_name text NOT NULL,
    age integer,
    is_active boolean DEFAULT TRUE,
    country text NOT NULL,
    personal_best interval,
    season_best interval,
    CONSTRAINT runners_pk PRIMARY KEY (id) -- definimos la clave principal de la tabla
);

CREATE INDEX runners_country
ON runners (country); -- creamos un índice en la columna country para optimizar las consultas que filtren por país

CREATE INDEX runners_season_best
ON runners (season_best);

-- results
CREATE TABLE results (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(), -- generamos un UUID basado en la dirección MAC del servidor y la fecha/hora actual
    runner_id uuid NOT NULL,
    race_result interval NOT NULL,
    location text NOT NULL,
    position integer,
    year integer NOT NULL,
    CONSTRAINT results_pk PRIMARY KEY (id), -- definimos la clave principal de la tabla
    CONSTRAINT fk_results_runner_id FOREIGN KEY (runner_id) -- definimos una foreign key que referencia a la tabla runners. La columna runner_id de results referencia a la columna id de runners. Cuando se actualiza o elimina un registro en runners, no se realiza ninguna acción en results
        REFERENCES runners (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

