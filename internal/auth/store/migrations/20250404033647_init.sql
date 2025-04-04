-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial not null primary key,
    uuid uuid DEFAULT gen_random_uuid() INDEX idx_uuid,
    email varchar not null unique INDEX idx_email,
    encrypted_password varchar not null INDEX idx_encrypted_password,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null default false
);

CREATE TABLE users_access (
    id serial not null primary key,
    user_id int not null,
    service varchar not null unique,
    roles varchar not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    deleted boolean default false,
    INDEX idx_user_service (user_id, service) 
);


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
DROP TABLE users;
DROP TABLE users_access;
-- +goose StatementEnd
