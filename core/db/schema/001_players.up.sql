-- The schema for bank ledger data
CREATE SCHEMA IF NOT EXISTS cncraft;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cncraft.players
(
    iso_code               CHAR(3)                       NOT NULL PRIMARY KEY,
    name                   VARCHAR(128)                  NOT NULL,
    type                   VARCHAR(32)                   NOT NULL,
    decimal_places         INTEGER                       NOT NULL,
    decimal_precision      INTEGER                       NOT NULL,
    created_timestamp      TIMESTAMP WITHOUT TIME ZONE   NOT NULL
);
