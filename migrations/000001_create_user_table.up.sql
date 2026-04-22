CREATE TABLE users (
                       uuid UUID PRIMARY KEY not null ,
                       name TEXT NOT NULL,
                       email TEXT NOT NULL unique,
                       password TEXT not null
);