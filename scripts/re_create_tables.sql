-- DROP DATABASE IF EXISTS admin_section_go;
-- CREATE DATABASE admin_section_go;
-- GRANT ALL PRIVILEGES ON DATABASE admin_section_go to admin_section_go;

CREATE TABLE IF NOT EXISTS users (
  id         SERIAL,
  firstName  TEXT NOT NULL,
  lastName   TEXT NOT NULL,
  middleName TEXT NOT NULL,
  email      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
  id       SERIAL,
  username VARCHAR(16) UNIQUE,
  password CHAR(64)
);

CREATE TABLE IF NOT EXISTS workerTypes (
  id             SERIAL UNIQUE,
  workerTypeName VARCHAR(128) UNIQUE
);

CREATE TABLE IF NOT EXISTS workers (
  id         SERIAL,
  name       VARCHAR(256),
  workerType INTEGER REFERENCES workerTypes (id)
);