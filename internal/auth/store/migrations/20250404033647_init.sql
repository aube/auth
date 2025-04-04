-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial not null primary key,
    uuid uuid DEFAULT gen_random_uuid() not null unique,
    email varchar not null unique,
    encrypted_password varchar not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

create INDEX idx_uuid on users (uuid);

create INDEX idx_email on users (email);

create INDEX idx_encrypted_password on users (encrypted_password);


CREATE TABLE users_access (
    id serial not null primary key,
    user_id int not null,
    service varchar not null unique,
    roles varchar not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean default false
);

create INDEX idx_user_service on users_access (user_id, service);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER users_access_updated_at_trigger
BEFORE UPDATE ON users_access
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER users_updated_at_trigger ON users;
DROP TRIGGER users_access_updated_at_trigger ON users_access;
DROP FUNCTION update_updated_at();
drop INDEX idx_encrypted_password;
drop INDEX idx_email;
drop INDEX idx_uuid;
drop INDEX idx_user_service;
DROP TABLE users;
DROP TABLE users_access;

-- +goose StatementEnd
