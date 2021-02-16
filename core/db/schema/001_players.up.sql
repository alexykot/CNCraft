-- The schema for bank ledger data
CREATE SCHEMA IF NOT EXISTS cncraft;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cncraft.players
(
    id         UUID                        NOT NULL PRIMARY KEY,
    username   VARCHAR(128)                NOT NULL,
    position_x DOUBLE PRECISION            NOT NULL,
    position_y DOUBLE PRECISION            NOT NULL,
    position_z DOUBLE PRECISION            NOT NULL,
    yaw        FLOAT                       NOT NULL,
    pitch      FLOAT                       NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
