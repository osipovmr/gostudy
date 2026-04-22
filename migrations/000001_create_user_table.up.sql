CREATE TABLE "user" (
                       id UUID PRIMARY KEY not null ,
                       name TEXT NOT NULL,
                       email TEXT NOT NULL unique 
);