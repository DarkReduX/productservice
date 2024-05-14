CREATE TABLE IF NOT EXISTS products
(
    id          uuid         NOT NULL,
    user_id     uuid         NOT NULL,
    name        VARCHAR(32)  NOT NULL,
    description varchar(600) NOT NULL,
    price       INTEGER CHECK (price > 0),
    PRIMARY KEY (id, user_id)
)