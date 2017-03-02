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
    disabled_until TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    last_login TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    UNIQUE(email, display_name)
);

insert into users(email, hash, display_name) VALUES ('trevor@youfie.io', '$2a$10$BW7ryN6m2lM8z7f57H69a.CXdQozmgrml20tf82lDE193nAozqKEa', 'trevor');
