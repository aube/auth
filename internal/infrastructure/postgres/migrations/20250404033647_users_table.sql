-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id serial not null primary key,
    username varchar not null unique,
    email varchar not null unique,
    phone integer DEFAULT null,
    encrypted_password varchar not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted boolean not null DEFAULT false
);

create INDEX users_username on users (username);

create INDEX users_encrypted_password on users (encrypted_password);


CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER users_updated_at_trigger ON users;

drop INDEX users_encrypted_password;
drop INDEX users_username;

DROP TABLE users;

-- +goose StatementEnd
