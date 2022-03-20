DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'EVENT_STATUS') THEN
        CREATE TYPE EVENT_STATUS AS ENUM ('available', 'finished', 'canceled', 'sold_out');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'EVENT_AGE_GROUP') THEN
        CREATE TYPE EVENT_AGE_GROUP AS ENUM ('0', '10', '12', '14', '16', '18');
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS "events" (
    id UUID UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    description TEXT DEFAULT NULL,
    organizer_account_id UUID NOT NULL,
    age_group EVENT_AGE_GROUP NOT NULL,
    maximum_capacity INTEGER NOT NULL CHECK (maximum_capacity >= 1),
    status EVENT_STATUS NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    duration SMALLINT NOT NULL CHECK (duration >= 1),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (organizer_account_id) REFERENCES accounts(id) ON DELETE SET NULL
);