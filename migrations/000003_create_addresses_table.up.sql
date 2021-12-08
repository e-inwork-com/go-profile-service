CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    owner UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE,
    street text NOT NULL,
    post_code text NOT NULL,
    city text NOT NULL,
    country_code text NOT NULL,
    version integer NOT NULL DEFAULT 1
);