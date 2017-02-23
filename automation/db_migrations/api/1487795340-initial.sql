DROP SCHEMA PUBLIC CASCADE;

CREATE SCHEMA PUBLIC;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    admin BOOLEAN DEFAULT false,

    disabled BOOLEAN DEFAULT false,
    disabled_until TIMESTAMP DEFAULT current_timestamp,

    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,

    UNIQUE(email, display_name)
);

