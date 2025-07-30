-- +goose Up
-- +goose StatementBegin

CREATE TABLE node_page (
    id SERIAL not null primary key,
    node_id smallint NOT NULL,
    page_id smallint NOT NULL,
    sort smallint not null default '0',
    pinned smallint not null default '0'
);


CREATE INDEX node_page_node_id on node_page (node_id);
CREATE INDEX node_page_page_id on node_page (page_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX node_page_node_id;
DROP INDEX node_page_page_id;

DROP TABLE node_page;

-- +goose StatementEnd