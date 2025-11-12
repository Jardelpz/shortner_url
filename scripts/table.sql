CREATE TABLE IF NOT EXISTS table_url (
    id serial PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL
)