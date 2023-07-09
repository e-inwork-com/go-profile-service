CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at_dt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    profile_user_s UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE,
    profile_name_t char varying(255) NOT NULL,
    profile_picture_s char varying(512),
    version integer NOT NULL DEFAULT 1
);