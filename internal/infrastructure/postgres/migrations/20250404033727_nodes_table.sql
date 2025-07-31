-- +goose Up
-- +goose StatementBegin


CREATE TABLE nodes (
    id SERIAL not null primary key,
    page_id INTEGER NOT NULL,

    name varchar(256) NOT NULL,
    menu_title varchar(256) NOT NULL default '',
    path varchar(512) NOT NULL default '/',
    
    parent INTEGER NOT NULL DEFAULT '0',
    sort smallint not null default '0',
    level smallint not null default '0',
    children smallint not null default '0',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    published boolean not null default false,
    deleted boolean not null default false
);


CREATE INDEX nodes_page_id on nodes (page_id);

CREATE TRIGGER nodes_updated_at_trigger
BEFORE UPDATE ON nodes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER nodes_updated_at_trigger ON nodes;

DROP INDEX nodes_page_id;

DROP TABLE nodes;

-- +goose StatementEnd