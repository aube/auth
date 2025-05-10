-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id serial not null primary key,
    uuid uuid DEFAULT gen_random_uuid() not null unique,
    email varchar not null,
    encrypted_password varchar not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

CREATE INDEX idx_uuid on users (uuid);
CREATE INDEX idx_email on users (email);

CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER users_updated_at_trigger ON users;

DROP INDEX idx_email;
DROP INDEX idx_uuid;

DROP TABLE users;

-- +goose StatementEnd
