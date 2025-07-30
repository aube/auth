-- +goose Up
-- +goose StatementBegin

CREATE TABLE images (
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

CREATE INDEX images_user_id on images (user_id);
CREATE INDEX images_name on images (name);
CREATE INDEX images_category on images (category);

CREATE TRIGGER images_updated_at_trigger
BEFORE UPDATE ON images
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER images_updated_at_trigger ON images;

DROP INDEX images_user_id;
DROP INDEX images_name;
DROP INDEX images_category;

DROP TABLE images;

-- +goose StatementEnd
