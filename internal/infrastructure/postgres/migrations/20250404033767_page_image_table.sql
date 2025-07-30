-- +goose Up
-- +goose StatementBegin

CREATE TABLE page_image (
    id SERIAL not null primary key,
    node_id smallint NOT NULL,
    page_id smallint NOT NULL,
    sort smallint not null default '0',
    pinned smallint not null default '0'
);


CREATE INDEX page_image_node_id on page_image (node_id);
CREATE INDEX page_image_page_id on page_image (page_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX page_image_node_id;
DROP INDEX page_image_page_id;

DROP TABLE page_image;

-- +goose StatementEnd