CREATE TABLE cocktails (
    id SERIAL PRIMARY KEY,
    name varchar(256) UNIQUE,
    ingredients text[],
    instructions text[]
)

