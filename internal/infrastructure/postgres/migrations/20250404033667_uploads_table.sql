-- +goose Up
-- +goose StatementBegin

CREATE TABLE uploads (
    id serial not null primary key,
    user_id bigint not null,
    uuid uuid not null,
    name varchar not null,
    category varchar default '',
    size bigint default 0,
    content_type varchar not null,
    description text default '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null DEFAULT false
);

CREATE INDEX uploads_user_id on uploads (user_id);
CREATE INDEX uploads_name on uploads (name);
CREATE INDEX uploads_category on uploads (category);

CREATE TRIGGER uploads_updated_at_trigger
BEFORE UPDATE ON uploads
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER uploads_updated_at_trigger ON uploads;

DROP INDEX uploads_user_id;
DROP INDEX uploads_name;
DROP INDEX uploads_category;

DROP TABLE uploads;

-- +goose StatementEnd
