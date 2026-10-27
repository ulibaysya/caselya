-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    image_url TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS item_tag_keys (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    is_public BOOL NOT NULL DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS item_tag_values (
    id BIGSERIAL PRIMARY KEY,
    tag_key_id BIGINT NOT NULL REFERENCES item_tag_keys(id),
    name TEXT NOT NULL,
    UNIQUE (tag_key_id, name)
);
CREATE TABLE IF NOT EXISTS item_tag_relations (
    item_id BIGINT NOT NULL REFERENCES items(id),
    tag_value_id BIGINT NOT NULL REFERENCES item_tag_values(id),
    PRIMARY KEY (item_id, tag_value_id)
);

CREATE TABLE IF NOT EXISTS items_price_providers (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS items_prices (
    item_id BIGINT NOT NULL REFERENCES items(id),
    price BIGINT NOT NULL CHECK (price >= 0),
    provider_id BIGINT NOT NULL REFERENCES items_price_providers(id),
    PRIMARY KEY (item_id, provider_id)
);

-- CREATE TABLE IF NOT EXISTS containers (
--     id BIGSERIAL PRIMARY KEY,
--     price BIGINT NOT NULL CHECK (price >= 0)
-- );
-- CREATE TABLE IF NOT EXISTS containers_items (
--     container_id BIGINT NOT NULL REFERENCES containers(id),
--     item_id  BIGINT NOT NULL REFERENCES items(id),
--     weight BIGINT NOT NULL CHECK (weight > 0),
--     PRIMARY KEY (container_id, item_id)
-- );
--
-- CREATE TABLE IF NOT EXISTS users (
--     id BIGSERIAL PRIMARY KEY,
--     balance BIGINT NOT NULL CHECK (balance >= 0)
-- );
-- CREATE TABLE IF NOT EXISTS users_items (
--     id BIGSERIAL PRIMARY KEY,
--     user_id BIGINT NOT NULL REFERENCES users(id),
--     item_id BIGINT NOT NULL REFERENCES items(id)
-- );

-- +goose Down
-- DROP TABLE users_items;
-- DROP TABLE users;
--
-- DROP TABLE containers_items;
-- DROP TABLE containers;

DROP TABLE items_prices;
DROP TABLE items_price_providers;

DROP TABLE item_tag_relations;
DROP TABLE item_tag_values;
DROP TABLE item_tag_keys;

DROP TABLE items;
