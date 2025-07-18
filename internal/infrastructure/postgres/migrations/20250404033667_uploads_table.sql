-- +goose Up
-- +goose StatementBegin

CREATE TABLE uploads (
    id serial not null primary key,
    user_id integer not null,
    uuid uuid not null,
    size bigint default 0,
    type varchar not null,
    name varchar not null,
    description text default '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null DEFAULT false
);

CREATE INDEX uploads_user_id on uploads (user_id);

CREATE TRIGGER uploads_updated_at_trigger
BEFORE UPDATE ON uploads
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER uploads_updated_at_trigger ON uploads;

DROP INDEX uploads_user_id;

DROP TABLE uploads;

-- +goose StatementEnd
