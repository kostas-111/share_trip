-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS trip_domain;

CREATE TYPE trip_domain.trip_status AS ENUM (
    'draft',
    'published',
    'canceled',
    'completed'
);

CREATE TABLE trip_domain.trips (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL,
    from_point TEXT NOT NULL,
    to_point TEXT NOT NULL,
    departure_time TIMESTAMPTZ NOT NULL,
    seats INT NOT NULL CHECK (seats > 0),
    status trip_domain.trip_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trip_domain.trip_history (
    id UUID PRIMARY KEY,
    trip_id UUID NOT NULL REFERENCES trip_domain.trips(id) ON DELETE CASCADE,
    from_status trip_domain.trip_status,
    to_status trip_domain.trip_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trip_domain.trip_history;
DROP TABLE IF EXISTS trip_domain.trips;
DROP TYPE IF EXISTS trip_domain.trip_status;
DROP SCHEMA IF EXISTS trip_domain.trip_domain;
-- +goose StatementEnd
