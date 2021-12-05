CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    owner UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE,
    first_name text NOT NULL,
    last_name text NOT NULL,
    version integer NOT NULL DEFAULT 1
);