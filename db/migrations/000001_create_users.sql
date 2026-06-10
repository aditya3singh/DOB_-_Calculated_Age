-- 000001_create_users.sql
-- Migration: create the users table

CREATE TABLE IF NOT EXISTS users (
    id   SERIAL      PRIMARY KEY,
    name TEXT        NOT NULL,
    dob  DATE        NOT NULL
);
