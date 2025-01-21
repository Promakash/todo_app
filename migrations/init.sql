CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(1024),
    is_done BOOLEAN
)