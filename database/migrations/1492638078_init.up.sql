CREATE TABLE cocktails (
    id SERIAL PRIMARY KEY NOT NULL,
    name varchar(256) UNIQUE,
    ingredients jsonb,
    instructions jsonb
)
