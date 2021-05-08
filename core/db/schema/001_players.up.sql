CREATE SCHEMA IF NOT EXISTS cncraft;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cncraft.players
(
    id             UUID                        NOT NULL,
    conn_id        UUID                        NULL,
    username       VARCHAR(128)                NOT NULL,

    position_x     DOUBLE PRECISION            NOT NULL,
    position_y     DOUBLE PRECISION            NOT NULL,
    position_z     DOUBLE PRECISION            NOT NULL,
    yaw            FLOAT                       NOT NULL,
    pitch          FLOAT                       NOT NULL,
    on_ground      BOOL                        NOT NULL DEFAULT TRUE,

    current_hotbar SMALLINT                    NOT NULL DEFAULT 0,

    created_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ix_players_username ON cncraft.players (username);

CREATE TABLE cncraft.inventory
(
    player_id   UUID REFERENCES cncraft.players (id) ON DELETE CASCADE,
    slot_number SMALLINT NOT NULL DEFAULT 0,
    item_id     SMALLINT NOT NULL DEFAULT 0,
    item_count  SMALLINT NOT NULL DEFAULT 0,

    PRIMARY KEY (player_id, slot_number)
);
