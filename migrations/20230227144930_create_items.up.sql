CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;

CREATE TABLE books(
  id serial PRIMARY KEY NOT NULL,
  title character varying NOT NULL,
  genre character varying NOT NULL,
  author character varying NOT NULL,
  price float
)
