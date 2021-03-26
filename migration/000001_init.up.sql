CREATE TABLE users (
     id uuid NOT NULL,
     name varchar(32) NOT NULL UNIQUE,
     password varchar(100),
     PRIMARY KEY (id)
);