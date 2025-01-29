-- +migrate Up
CREATE TABLE tasks(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(1024),
    is_done BOOLEAN
);

-- +migrate Down
DROP TABLE IF EXISTS tasks;