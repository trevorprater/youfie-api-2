DROP SCHEMA PUBLIC CASCADE;

CREATE SCHEMA PUBLIC;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    hash TEXT NOT NULL,
    display_name TEXT NOT NULL UNIQUE,
    admin BOOLEAN DEFAULT false,

    disabled BOOLEAN DEFAULT false,
    disabled_until TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

    last_login TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp

);

CREATE TABLE photos(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID REFERENCES users (id),
    format TEXT NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    storage_url TEXT,
    latitude NUMERIC(10, 6),
    longitude NUMERIC(10, 6),
    processed BOOLEAN DEFAULT false,
    processing BOOLEAN default false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE faces(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    photo_id UUID NOT NULL REFERENCES photos(id),
    feature_vector TEXT NOT NULL,

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
    is_match BOOLEAN NOT NULL DEFAULT true,
    user_acknowledged BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE conversations(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    photo_id UUID NOT NULL REFERENCES photos (id),
    owner_id UUID NOT NULL REFERENCES users (id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE conversation_participants(
    conversation_id UUID PRIMARY KEY REFERENCES conversations (id),
    user_id UUID NOT NULL REFERENCES users (id),
    user_approved BOOLEAN DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE conversation_messages(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations (id),
    owner_id UUID NOT NULL REFERENCES users (id),
    message_text TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);
