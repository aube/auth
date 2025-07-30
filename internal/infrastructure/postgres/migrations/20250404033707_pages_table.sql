-- +goose Up
-- +goose StatementBegin


CREATE TABLE pages (
    id SERIAL not null primary key,
    name varchar(256) NOT NULL,

    meta text NOT NULL default '',
    title varchar(1024) NOT NULL default '',
    category varchar NOT NULL default '',
    template varchar NOT NULL default '',

    h1 varchar(1024) NOT NULL default '',
    content text NOT NULL default '',
    content_short varchar(4096) NOT NULL default '',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- show_children boolean not null default false,
    published boolean not null default false,
    deleted boolean not null default false
);


CREATE TRIGGER pages_updated_at_trigger
BEFORE UPDATE ON pages
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER pages_updated_at_trigger ON pages;

DROP TABLE pages;


-- +goose StatementEnd