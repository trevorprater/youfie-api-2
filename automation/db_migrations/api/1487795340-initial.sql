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

CREATE TABLE photos(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID REFERENCES users (id),
    format TEXT NOT NULL,
    content BYTEA NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    storage_url TEXT,
    latitude NUMERIC(10, 6),
    longitude NUMERIC(10, 6),
    processed BOOLEAN DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE faces(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    photo_id UUID NOT NULL REFERENCES photos(id),
    feature_vector DECIMAL[128] NOT NULL,

    bb_top_left_x INT NOT NULL,
    bb_top_left_y INT NOT NULL,
    bb_top_right_x INT NOT NULL,
    bb_top_right_y INT NOT NULL,
    bb_bottom_left_x INT NOT NULL,
    bb_bottom_left_y INT NOT NULL,
    bb_bottom_right_x INT NOT NULL,
    bb_bottom_right_y INT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE matches(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    photo_id UUID NOT NULL REFERENCES photos (id),
    face_id UUID NOT NULL REFERENCES faces (id),
    user_id UUID NOT NULL REFERENCES users (id),
    confidence DECIMAL,
    confirmed BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

