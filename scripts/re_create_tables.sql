-- DROP DATABASE IF EXISTS admin_section_go;
-- CREATE DATABASE admin_section_go;
-- GRANT ALL PRIVILEGES ON DATABASE admin_section_go to admin_section_go;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL,
  firstName TEXT NOT NULL,
  lastName TEXT NOT NULL,
  middleName TEXT NOT NULL,
  email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
  id SERIAL,
  username VARCHAR(16) UNIQUE,
  password CHAR(64)
);
