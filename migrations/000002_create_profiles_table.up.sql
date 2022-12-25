CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    profile_user UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE,
    profile_name char varying(100) NOT NULL,
    profile_picture char varying(512),
    version integer NOT NULL DEFAULT 1
);