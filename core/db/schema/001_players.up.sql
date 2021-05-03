CREATE SCHEMA IF NOT EXISTS cncraft;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cncraft.players
(
    id           UUID                        NOT NULL PRIMARY KEY,
    username     VARCHAR(128)                NOT NULL,

    position_x   DOUBLE PRECISION            NOT NULL,
    position_y   DOUBLE PRECISION            NOT NULL,
    position_z   DOUBLE PRECISION            NOT NULL,
    yaw          FLOAT                       NOT NULL,
    pitch        FLOAT                       NOT NULL,
    on_ground    BOOL                        NOT NULL DEFAULT TRUE,

    current_slot INT                         NOT NULL DEFAULT 0,
    slot0        INT                         NOT NULL DEFAULT 0,
    slot1        INT                         NOT NULL DEFAULT 0,
    slot2        INT                         NOT NULL DEFAULT 0,
    slot3        INT                         NOT NULL DEFAULT 0,
    slot4        INT                         NOT NULL DEFAULT 0,
    slot5        INT                         NOT NULL DEFAULT 0,
    slot6        INT                         NOT NULL DEFAULT 0,
    slot7        INT                         NOT NULL DEFAULT 0,

    created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX ix_players_username ON cncraft.players (username);
