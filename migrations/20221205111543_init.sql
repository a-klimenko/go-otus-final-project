-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS banners (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS slots (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS groups (
    id          UUID PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS rotations (
    id        UUID PRIMARY KEY,
    banner_id UUID,
    slot_id   UUID,
    group_id  UUID,
    shows     INT DEFAULT 1,
    clicks    INT DEFAULT 0,
    deleted_at TIMESTAMP DEFAULT NULL,

    FOREIGN KEY (banner_id) REFERENCES banners(id),
    FOREIGN KEY (slot_id)   REFERENCES slots(id),
    FOREIGN KEY (group_id)  REFERENCES groups(id),
    UNIQUE(banner_id, slot_id, group_id)
);

INSERT INTO slots (id, description)
VALUES (uuid_generate_v4(), 'left slot'), (uuid_generate_v4(), 'right slot'), (uuid_generate_v4(), 'footer slot');

INSERT INTO banners (id, description)
VALUES (uuid_generate_v4(), 'warm socks'), (uuid_generate_v4(), 'real estate advertising'), (uuid_generate_v4(), 'electronics store'),
       (uuid_generate_v4(), 'medicine advertising'), (uuid_generate_v4(), 'musical concert');

INSERT INTO groups (id, description)
VALUES (uuid_generate_v4(), 'aged'), (uuid_generate_v4(), 'youth'), (uuid_generate_v4(), 'man'), (uuid_generate_v4(), 'woman');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rotations;
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS slots;
DROP TABLE IF EXISTS groups;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
